package templates

const IndexCSS = `@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --primary: #007acc; /* VSCode Blue */
    --primary-foreground: #ffffff;
    --background: #1e1e1e; /* VSCode Dark */
    --foreground: #d4d4d4;
    --card: #252526;
    --card-foreground: #d4d4d4;
    --border: #3e3e42;
  }

  /* Light Theme Override */
  .light {
    --primary: #007acc;
    --primary-foreground: #ffffff;
    --background: #ffffff;
    --foreground: #000000;
    --card: #f3f3f3;
    --card-foreground: #000000;
    --border: #e0e0e0;
  }

  body {
    @apply bg-background text-foreground font-sans antialiased;
  }
}`

const AppCSS = `/* Additional custom styles if needed, mostly handled by Tailwind now */
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.react:hover {
  filter: drop-shadow(0 0 2em #61dafbaa);
}

@keyframes logo-spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@media (prefers-reduced-motion: no-preference) {
  a:nth-of-type(2) .logo {
    animation: logo-spin infinite 20s linear;
  }
}
`
