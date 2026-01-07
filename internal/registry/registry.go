package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type WodgeApp struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Port      int       `json:"port"`
	PID       int       `json:"pid"`
	StartTime time.Time `json:"start_time"`
}

type Registry struct {
	Apps map[string]WodgeApp `json:"apps"`
	mu   sync.Mutex
}

func getRegistryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".wodge")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "registry.json"), nil
}

func Load() (*Registry, error) {
	path, err := getRegistryPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Registry{Apps: make(map[string]WodgeApp)}, nil
	}
	if err != nil {
		return nil, err
	}

	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	if reg.Apps == nil {
		reg.Apps = make(map[string]WodgeApp)
	}
	return &reg, nil
}

func (r *Registry) Save() error {
	path, err := getRegistryPath()
	if err != nil {
		return err
	}

	// Clean up stale PIDs before saving
	for name, app := range r.Apps {
		if !isProcessRunning(app.PID) {
			delete(r.Apps, name)
		}
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (r *Registry) Register(name string, port int, projectPath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Reload first to get latest state
	current, err := Load()
	if err == nil {
		r.Apps = current.Apps
	}

	r.Apps[name] = WodgeApp{
		Name:      name,
		Path:      projectPath,
		Port:      port, // Go backend port
		PID:       os.Getpid(),
		StartTime: time.Now(),
	}
	return r.Save()
}

func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, err := Load()
	if err == nil {
		r.Apps = current.Apps
	}

	delete(r.Apps, name)
	return r.Save()
}

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 checks if process exists
	err = process.Signal(syscall.Signal(0)) // syscall.Signal(0) is portable enough usually
	return err == nil
}

func GetFreePort(startPort int) int {
	port := startPort
	for {
		if !isPortInUse(port) {
			return port
		}
		port++
		if port > 65535 {
			return 0 // No free ports?!
		}
	}
}

func isPortInUse(port int) bool {
	// Simple check by trying to listen
	// NOT perfect due to race conditions but good enough for dev
	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], []byte{0, 0, 0, 0}) // 0.0.0.0

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return true // Assume used if error
	}
	defer syscall.Close(fd)

	// Set SO_REUSEADDR so we don't get stuck waiting
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return true
	}

	if err := syscall.Bind(fd, &addr); err != nil {
		return true // Port is used
	}

	return false
}
