package templates

const ComponentButton = `import React from 'react';
import { motion } from 'framer-motion';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', ...props }, ref) => {
    const variants = {
      primary: 'bg-primary text-primary-foreground hover:opacity-90',
      secondary: 'bg-card text-card-foreground hover:bg-card/80',
      outline: 'border border-border bg-transparent hover:bg-card',
      ghost: 'bg-transparent hover:bg-card text-foreground',
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
import { motion } from 'framer-motion';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const Card = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
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
