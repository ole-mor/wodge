package templates

const TailwindConfig = `import type { Config } from 'tailwindcss';

export default {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    fontFamily: {
      sans: ['Assistant', 'sans-serif'],
    },
    extend: {
      colors: {
        primary: 'var(--primary)',
        'primary-foreground': 'var(--primary-foreground)',
        background: 'var(--background)',
        foreground: 'var(--foreground)',
        card: 'var(--card)',
        'card-foreground': 'var(--card-foreground)',
        border: 'var(--border)',
      },
    },
  },
  plugins: [],
} satisfies Config;`

const PostCSSConfig = `export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}`
