package main

import (
	"encoding/json"
	"net/http"
)

// ApiResponse represents a standard API response
type ApiResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Status: status,
		Data:   data,
	})
}

// WriteError writes an error JSON response
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Status: status,
		Error:  message,
	})
}

// Handler type for API endpoint handlers
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// MethodRouter routes HTTP methods to their respective handlers
type MethodRouter struct {
	get    HandlerFunc
	post   HandlerFunc
	put    HandlerFunc
	delete HandlerFunc
}

// NewMethodRouter creates a new method router
func NewMethodRouter() *MethodRouter {
	return &MethodRouter{}
}

// GET registers a GET handler
func (mr *MethodRouter) GET(fn HandlerFunc) *MethodRouter {
	mr.get = fn
	return mr
}

// POST registers a POST handler
func (mr *MethodRouter) POST(fn HandlerFunc) *MethodRouter {
	mr.post = fn
	return mr
}

// PUT registers a PUT handler
func (mr *MethodRouter) PUT(fn HandlerFunc) *MethodRouter {
	mr.put = fn
	return mr
}

// DELETE registers a DELETE handler
func (mr *MethodRouter) DELETE(fn HandlerFunc) *MethodRouter {
	mr.delete = fn
	return mr
}

// Handler returns an http.HandlerFunc that routes based on method
func (mr *MethodRouter) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if mr.get != nil {
				mr.get(w, r)
			} else {
				WriteError(w, http.StatusMethodNotAllowed, "GET not allowed")
			}
		case http.MethodPost:
			if mr.post != nil {
				mr.post(w, r)
			} else {
				WriteError(w, http.StatusMethodNotAllowed, "POST not allowed")
			}
		case http.MethodPut:
			if mr.put != nil {
				mr.put(w, r)
			} else {
				WriteError(w, http.StatusMethodNotAllowed, "PUT not allowed")
			}
		case http.MethodDelete:
			if mr.delete != nil {
				mr.delete(w, r)
			} else {
				WriteError(w, http.StatusMethodNotAllowed, "DELETE not allowed")
			}
		default:
			WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}
