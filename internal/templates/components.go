package templates

const ComponentButton = `import React from 'react';
import { motion, type HTMLMotionProps } from 'framer-motion';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

interface ButtonProps extends HTMLMotionProps<'button'> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive';
  size?: 'sm' | 'md' | 'lg';
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', ...props }, ref) => {
    const variants = {
      primary: 'bg-primary text-primary-foreground hover:opacity-90',
      secondary: 'bg-card text-card-foreground hover:bg-card/80',
      outline: 'border border-border bg-transparent hover:bg-card',
      ghost: 'bg-transparent hover:bg-card text-foreground',
      destructive: 'bg-destructive text-destructive-foreground hover:opacity-90',
    };

    const sizes = {
      sm: 'h-8 px-3 text-xs',
      md: 'h-10 px-4 py-2',
      lg: 'h-12 px-8 text-lg',
    };

    return (
      <motion.button
        ref={ref}
        whileTap={{ scale: 0.95 }}
        whileHover={{ scale: 1.02 }}
        className={cn(
          'inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary disabled:pointer-events-none disabled:opacity-50',
          variants[variant],
          sizes[size],
          className
        )}
        {...props}
      />
    );
  }
);
Button.displayName = 'Button';
`

const ComponentCard = `import React from 'react';
import { motion, HTMLMotionProps } from 'framer-motion';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const Card = React.forwardRef<HTMLDivElement, HTMLMotionProps<'div'>>(
  ({ className, ...props }, ref) => (
    <motion.div
      ref={ref}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className={cn(
        'rounded-lg border border-border bg-card text-card-foreground shadow-sm',
        className
      )}
      {...props}
    />
  )
);
Card.displayName = 'Card';

export const CardHeader = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cn('flex flex-col space-y-1.5 p-6', className)}
      {...props}
    />
  )
);
CardHeader.displayName = 'CardHeader';

export const CardTitle = React.forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLHeadingElement>>(
  ({ className, ...props }, ref) => (
    <h3
      ref={ref}
      className={cn('text-2xl font-semibold leading-none tracking-tight', className)}
      {...props}
    />
  )
);
CardTitle.displayName = 'CardTitle';

export const CardContent = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div ref={ref} className={cn('p-6 pt-0', className)} {...props} />
  )
);
CardContent.displayName = 'CardContent';
`

const ComponentInput = `import React from 'react';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          'flex h-10 w-full rounded-md border border-border bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary disabled:cursor-not-allowed disabled:opacity-50',
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);
Input.displayName = 'Input';
`

const ComponentNavbar = `import React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { useTheme } from '@/context/ThemeProvider';
import LightModeIcon from '@mui/icons-material/LightMode';
import DarkModeIcon from '@mui/icons-material/DarkMode';
import DashboardIcon from '@mui/icons-material/Dashboard';

export function Navbar() {
  const { theme, toggleTheme } = useTheme();

  return (
    <nav className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-50">
      <div className="flex h-16 items-center px-4 container mx-auto">
        <div className="mr-8 flex items-center space-x-2">
          <DashboardIcon className="text-primary" />
          <span className="text-lg font-bold">Wodge App</span>
        </div>
        
        <div className="flex items-center space-x-6 text-sm font-medium">
          <Link to="/" className="transition-colors hover:text-primary">Home</Link>
          <Link to="/about" className="transition-colors hover:text-primary">About</Link>
        </div>

        <div className="ml-auto flex items-center space-x-4">
          <Button variant="ghost" size="sm" onClick={toggleTheme}>
            {theme === 'dark' ? <LightModeIcon /> : <DarkModeIcon />}
          </Button>
        </div>
      </div>
    </nav>
  );
}
`

const ComponentQastTest = `import React, { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { qast } from '@/api/qast';

export function QastTest() {
  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');
  const [message, setMessage] = useState('');
  const [response, setResponse] = useState<any>(null);

  const testConnectivity = async () => {
    setStatus('loading');
    setMessage('Pinging Qast via Proxy...');
    setResponse(null);
    try {
        const res = await qast.ask("PING_CONNECTIVITY_TEST");
        console.log("Qast Response:", res);
        setResponse(res);
        setStatus('success');
        setMessage('Connected successfully to Qast Network!');
    } catch (e) {
        console.error("Connectivity Test Failed:", e);
        setStatus('error');
        setMessage(e instanceof Error ? e.message : 'Connection failed');
        setResponse(null);
    }
  };

  const testStreaming = async () => {
    setStatus('loading');
    setMessage('Testing Streaming Channel...');
    setResponse([]);
    try {
        await qast.chatStream("PING_STREAM_TEST", (event) => {
            console.log("Stream Event:", event);
            setResponse((prev: any) => {
                const arr = Array.isArray(prev) ? prev : [];
                return [...arr, event];
            });
            if (event.type === 'status') setMessage(event.data);
            if (event.type === 'done') {
                setStatus('success');
                setMessage('Streaming Test Completed!');
            }
        });
    } catch (e) {
        console.error("Stream Test Failed:", e);
        setStatus('error');
        setMessage('Stream Error: ' + (e instanceof Error ? e.message : String(e)));
    }
  };

  return (
    <Card className="w-full max-w-md mx-auto mt-8 border-2 border-border/50">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
            <span>Qast Connectivity</span>
            <div className="flex items-center gap-2">
                <span className="text-xs font-normal text-muted-foreground">Status:</span>
                {status === 'success' && <div className="h-3 w-3 rounded-full bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.6)]" />}
                {status === 'error' && <div className="h-3 w-3 rounded-full bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.6)]" />}
                {status === 'loading' && <div className="h-3 w-3 rounded-full bg-yellow-400 animate-pulse" />}
                {status === 'idle' && <div className="h-3 w-3 rounded-full bg-gray-600" />}
            </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="text-xs px-3 py-2 bg-muted/50 rounded border border-border/50 font-mono text-muted-foreground flex justify-between">
            <span>Proxy Target:</span>
            <span className="text-foreground">http://localhost:9988</span>
        </div>
        
        <div className="flex gap-2">
            <Button 
                onClick={testConnectivity} 
                disabled={status === 'loading'} 
                className="flex-1 font-bold"
                variant={status === 'error' ? 'destructive' : 'primary'}
            >
                Ping
            </Button>
            <Button 
                onClick={testStreaming} 
                disabled={status === 'loading'} 
                className="flex-1 font-bold"
                variant="outline"
            >
                Stream Test
            </Button>
        </div>

        {message && (
            <div className={"p-3 rounded-md text-sm font-medium border " + 
                (status === 'error' ? 'bg-destructive/10 border-destructive/20 text-destructive' : 
                status === 'success' ? 'bg-green-500/10 border-green-500/20 text-green-500' : 'bg-muted border-transparent')
            }>
                {message}
            </div>
        )}

        {response && (
            <div className="mt-4 space-y-2">
                <div className="text-xs font-semibold text-muted-foreground">Response Dump:</div>
                <div className="p-3 bg-black/40 rounded-md text-xs font-mono overflow-auto max-h-40 border border-border/30 shadow-inner">
                    <pre>{JSON.stringify(response, null, 2)}</pre>
                </div>
            </div>
        )}
      </CardContent>
    </Card>
  );
}
`

const ComponentTokenManager = `// A utility to manage Client-Side PII Tokens
// This mocks a secure storage. In production, this would be Encrypted LocalStorage.

const STORAGE_KEY = 'wodge_token_map';

export const TokenManager = {
  getMap(): Record<string, string> {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      return raw ? JSON.parse(raw) : {};
    } catch (e) {
      console.error("Failed to parse token map", e);
      return {};
    }
  },

  saveMap(newMap: Record<string, string>) {
    const current = this.getMap();
    const updated = { ...current, ...newMap };
    localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
  },

  get(token: string): string | null {
    const map = this.getMap();
    return map[token] || null;
  },

  // Replaces all known tokens in text with their real values
  rehydrate(text: string): string {
    const map = this.getMap();
    let rehydrated = text;
    
    // Sort keys by length desc to avoid partial replacements if tokens overlap
    const tokens = Object.keys(map).sort((a, b) => b.length - a.length);
    
    for (const token of tokens) {
      // Global replace of the token
      // Escape generic regex characters in token if necessary (usually tokens are [TYPE_N])
      const escapedToken = token.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&');
      const regex = new RegExp(escapedToken, 'g');
      rehydrated = rehydrated.replace(regex, map[token]);
    }
    
    return rehydrated;
  }
};
`

const ComponentSecureChat = `import React, { useState, useRef, useEffect } from 'react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { qast } from '@/api/qast';
import { TokenManager } from '@/utils/TokenManager'; 
import { motion, AnimatePresence } from 'framer-motion';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string; // This is the displayed content (rehydrated)
  isSanitized?: boolean; // If true, it was processed securely
}

export function SecureChat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [loadingStatus, setLoadingStatus] = useState("");
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages, loadingStatus]);

  const handleSend = async () => {
    if (!input.trim() || isLoading) return;

    const userMsg: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: input,
    };

    setMessages(prev => [...prev, userMsg]);
    setInput('');
    setIsLoading(true);
    setLoadingStatus("Securely processing...");

    // Preliminary empty bot message to fill
    const botMsgId = (Date.now() + 1).toString();
    const botMsg: Message = {
      id: botMsgId,
      role: 'assistant',
      content: '', // Will stream in
      isSanitized: true
    };
    setMessages(prev => [...prev, botMsg]);

    try {
        let currentContent = "";

        await qast.chatStream(userMsg.content, (event) => {
            if (event.type === 'status') {
                setLoadingStatus(event.data);
            } else if (event.type === 'token_map') {
                TokenManager.saveMap(event.data);
                console.log("SecureChat: Updated Token Map", event.data);
            } else if (event.type === 'chunk') {
                // Accumulate and rehydrate on every chunk
                // console.log("SecureChat Chunk Rx:", event.data); // DEBUG: Uncomment if needed
                currentContent += event.data;
                const rehydrated = TokenManager.rehydrate(currentContent);
                
                setMessages(prev => prev.map(m => 
                    m.id === botMsgId ? { ...m, content: rehydrated } : m
                ));
            } else if (event.type === 'error') {
                console.error("Stream Error:", event.data);
                // Optionally show error in chat?
            }
        });

    } catch (e) {
      console.error("SecureChat Error:", e);
      const errorMsg: Message = {
        id: (Date.now() + 2).toString(),
        role: 'assistant',
        content: "Error: " + (e instanceof Error ? e.message : ''),
      };
      setMessages(prev => [...prev, errorMsg]);
    } finally {
      setIsLoading(false);
      setLoadingStatus("");
    }
  };

  return (
    <Card className="w-full max-w-2xl mx-auto h-[600px] flex flex-col border-2 border-border/50 shadow-xl">
      <CardHeader className="border-b border-border/50 bg-muted/20">
        <CardTitle className="flex items-center gap-2">
          <span>Secure Chat</span>
          <span className="text-xs font-normal px-2 py-0.5 rounded-full bg-green-500/10 text-green-500 border border-green-500/20">
            End-to-End Privacy
          </span>
        </CardTitle>
      </CardHeader>
      
      <CardContent className="flex-1 overflow-hidden p-0 flex flex-col">
        {/* Messages Port */}
        <div className="flex-1 overflow-y-auto p-4 space-y-4" ref={scrollRef}>
          <AnimatePresence initial={false}>
            {messages.map(msg => (
              <motion.div
                key={msg.id}
                initial={{ opacity: 0, y: 10, scale: 0.95 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                className={"flex w-full " + (msg.role === 'user' ? 'justify-end' : 'justify-start')}
              >
                <div className={"max-w-[80%] rounded-2xl px-4 py-3 text-sm " + 
                  (msg.role === 'user' 
                    ? 'bg-primary text-primary-foreground rounded-br-none' 
                    : 'bg-muted text-foreground rounded-bl-none border border-border/50')
                }>
                  {msg.content}
                   {/* Verification Badge for Assistant */}
                   {msg.role === 'assistant' && msg.isSanitized && (
                    <div className="mt-2 text-[10px] opacity-70 flex items-center gap-1">
                        <div className="w-1.5 h-1.5 bg-green-500 rounded-full" />
                        <span>PII Rehydrated locally</span>
                    </div>
                  )}
                </div>
              </motion.div>
            ))}
          </AnimatePresence>
          
          {isLoading && (
             <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="flex justify-start">
               <div className="bg-muted px-4 py-3 rounded-2xl rounded-bl-none text-xs text-muted-foreground flex items-center gap-2">
                 <span>{loadingStatus || "Securely processing"}</span>
                 <div className="flex gap-1">
                    <div className="w-1.5 h-1.5 bg-foreground/50 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                    <div className="w-1.5 h-1.5 bg-foreground/50 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
                    <div className="w-1.5 h-1.5 bg-foreground/50 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
                 </div>
               </div>
             </motion.div>
          )}
        </div>

        {/* Input Area */}
        <div className="p-4 bg-background border-t border-border/50 flex gap-2">
          <Input 
            value={input}
            onChange={e => setInput(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleSend()}
            placeholder={isLoading ? "Waiting for secure response..." : "Type a message safely..."}
            className="flex-1 transition-all"
            disabled={isLoading}
          />
          <Button 
            onClick={handleSend} 
            disabled={isLoading || !input.trim()}
            variant="primary"
          >
            Send
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
`

const ComponentAuthProvider = `import React, { createContext, useContext, useState, useEffect } from 'react';
import { auth } from '@/api/auth';

interface User {
  id: string;
  email: string;
  role: string;
  [key: string]: any;
}

interface AuthContextType {
  user: User | null;
  accessToken: string | null;
  login: (username: string, pass: string) => Promise<void>;
  register: (email: string, username: string, pass: string, confirmPass: string, first: string, last: string) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(localStorage.getItem('access_token'));
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initAuth = async () => {
      const storedToken = localStorage.getItem('access_token');
      if (storedToken) {
        try {
           const res = await auth.verify(storedToken);
           setUser(res.user);
           setAccessToken(storedToken); // redundant but safe
        } catch (e) {
           console.warn("Auth check failed, logging out:", e);
           logout();
        }
      }
      setIsLoading(false);
    };

    initAuth();
  }, []);

  const login = async (username: string, pass: string) => {
    const res = await auth.login(username, pass);
    if (res.access_token) {
        setAccessToken(res.access_token);
        setUser(res.user);
        
        localStorage.setItem('access_token', res.access_token);
        localStorage.setItem('refresh_token', res.refresh_token);
        localStorage.setItem('user_profile', JSON.stringify(res.user));
    }
  };

  const register = async (email: string, username: string, pass: string, confirmPass: string, first: string, last: string) => {
    await auth.register(email, username, pass, confirmPass, first, last);
    await login(username, pass);
  };

  const logout = async () => {
    try {
      if (accessToken) {
        const refreshToken = localStorage.getItem('refresh_token');
        await auth.logout(accessToken, refreshToken || ""); 
      }
    } catch (e) {
      console.warn("Logout error:", e);
    } finally {
      setAccessToken(null);
      setUser(null);
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user_profile');
    }
  };

  return (
    <AuthContext.Provider value={{ 
        user, 
        accessToken, 
        login, 
        register, 
        logout, 
        isAuthenticated: !!accessToken,
        isLoading 
    }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
`

const ComponentProtectedRoute = `import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/context/AuthProvider';

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return <div className="flex h-screen items-center justify-center p-8 bg-background text-muted-foreground animate-pulse">Loading auth...</div>;
  }

  if (!isAuthenticated) {
    // Redirect to login, but save the current location they were trying to go to
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
}
`

const ComponentLoginPage = `import React, { useState } from 'react';
import { useAuth } from '@/context/AuthProvider';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { motion, AnimatePresence } from 'framer-motion';

export default function LoginPage() {
  const { login, register } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const from = (location.state as any)?.from?.pathname || '/';

  const [isLogin, setIsLogin] = useState(true);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const [form, setForm] = useState({
    email: '',
    username: '',
    password: '',
    confirmPassword: '',
    firstName: '',
    lastName: ''
  });

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(form.username, form.password);
      navigate(from, { replace: true });
    } catch (err: any) {
      setError(err.message || err.toString() || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    if (form.password !== form.confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    setError('');
    setLoading(true);
    try {
      await register(form.email, form.username, form.password, form.confirmPassword, form.firstName, form.lastName);
      navigate(from, { replace: true });
    } catch (err: any) {
      setError(err.message || err.toString() || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md border-border/50 shadow-2xl">
        <CardHeader className="text-center pb-2">
            <motion.div 
                key={isLogin ? 'login' : 'register'}
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3 }}
            >
                <CardTitle className="text-3xl font-bold tracking-tight mb-2">
                    {isLogin ? 'Welcome Back' : 'Create Account'}
                </CardTitle>
                <p className="text-sm text-muted-foreground">
                    {isLogin ? 'Enter your credentials to access your account' : 'Sign up to get started with Wodge'}
                </p>
            </motion.div>
        </CardHeader>
        <CardContent>
          <form onSubmit={isLogin ? handleLogin : handleRegister} className="space-y-4">
            <AnimatePresence mode="popLayout">
                {!isLogin && (
                    <motion.div 
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="space-y-4"
                    >
                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <label className="text-xs font-medium">First Name</label>
                                <Input placeholder="John" required={!isLogin} value={form.firstName} onChange={e => setForm({ ...form, firstName: e.target.value })} />
                            </div>
                            <div className="space-y-2">
                                <label className="text-xs font-medium">Last Name</label>
                                <Input placeholder="Doe" required={!isLogin} value={form.lastName} onChange={e => setForm({ ...form, lastName: e.target.value })} />
                            </div>
                        </div>
                        <div className="space-y-2">
                            <label className="text-xs font-medium">Email</label>
                            <Input 
                                type="email" 
                                placeholder="name@example.com" 
                                required={!isLogin}
                                value={form.email}
                                onChange={e => setForm({ ...form, email: e.target.value })}
                            />
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>

            <div className="space-y-2">
              <label className="text-xs font-medium">Username</label>
              <Input 
                placeholder="jdoe" 
                required 
                value={form.username}
                onChange={e => setForm({ ...form, username: e.target.value })}
              />
            </div>
            
            <div className="space-y-2">
              <label className="text-xs font-medium">Password</label>
              <Input 
                type="password" 
                placeholder="••••••••" 
                required 
                value={form.password}
                onChange={e => setForm({ ...form, password: e.target.value })}
              />
            </div>

            {!isLogin && (
                <div className="space-y-2">
                    <label className="text-xs font-medium">Confirm Password</label>
                    <Input 
                        type="password" 
                        placeholder="••••••••" 
                        required={!isLogin}
                        value={form.confirmPassword}
                        onChange={e => setForm({ ...form, confirmPassword: e.target.value })}
                    />
                </div>
            )}

            {error && (
              <div className="p-3 rounded bg-destructive/10 text-destructive text-sm font-medium border border-destructive/20">
                {error}
              </div>
            )}

            <Button type="submit" className="w-full font-bold" disabled={loading} size="lg">
              {loading ? (
                  <span className="flex items-center gap-2">
                      <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                      Processing...
                  </span>
              ) : (isLogin ? 'Sign In' : 'Sign Up')}
            </Button>
          </form>

          <div className="mt-6 text-center text-sm">
            <span className="text-muted-foreground">
                {isLogin ? "Don't have an account? " : "Already have an account? "}
            </span>
            <button 
                type="button"
                onClick={() => { setIsLogin(!isLogin); setError(''); }}
                className="font-semibold text-primary hover:underline focus:outline-none"
            >
                {isLogin ? 'Sign up' : 'Sign in'}
            </button>
          </div>
        </CardContent>
      </Card>
      
      {/* Background decoration */}
      <div className="fixed inset-0 -z-10 pointer-events-none">
          <div className="absolute top-0 right-0 w-1/2 h-1/2 bg-primary/5 blur-[120px] rounded-full translate-x-1/2 -translate-y-1/2" />
          <div className="absolute bottom-0 left-0 w-1/3 h-1/3 bg-blue-500/5 blur-[100px] rounded-full -translate-x-1/2 translate-y-1/2" />
      </div>
    </div>
  );
}
`
