package templates

const ComponentButton = `import React from 'react';
import { motion } from 'framer-motion';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
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
    console.log("Starting connectivity test to Qast Proxy...");
    try {
        // We use a simple ingest call or ask call to verify connectivity
        // Using 'ask' with a ping message
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
        
        <Button 
            onClick={testConnectivity} 
            disabled={status === 'loading'} 
            className="w-full font-bold relative overflow-hidden group"
            variant={status === 'error' ? 'destructive' : 'primary'}
        >
            {status === 'loading' ? 'Establishing Secure Link...' : 'Test Connection'}
        </Button>

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
                <div className="text-xs font-semibold text-muted-foreground">Last Response Packet:</div>
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
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

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

    try {
      // 1. Send text to Qast Secure Chat (User -> Proxy -> Qast)
      // The backend will Anonymize -> Store Mapping (Temp/Return) -> Ask LLM -> Return
      // We expect { llm_response: string, token_map: Record<string, string> }
      
      const res = await qast.chat(input);
      
      // 2. Update Client-Side Token Manager with new mappings
      if (res.token_map) {
        TokenManager.saveMap(res.token_map);
        console.log("SecureChat: Updated Token Map", res.token_map);
      }

      // 3. Rehydrate the response (Replace tokens with real names)
      const realText = TokenManager.rehydrate(res.llm_response);

      const botMsg: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: realText,
        isSanitized: true
      };

      setMessages(prev => [...prev, botMsg]);

    } catch (e) {
      console.error("SecureChat Error:", e);
      const errorMsg: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: "Error: Failed to establish secure link. " + (e instanceof Error ? e.message : ''),
      };
      setMessages(prev => [...prev, errorMsg]);
    } finally {
      setIsLoading(false);
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
                 <div className="w-2 h-2 bg-foreground/50 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                 <div className="w-2 h-2 bg-foreground/50 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
                 <div className="w-2 h-2 bg-foreground/50 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
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
            placeholder="Type a message safely..."
            className="flex-1"
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
