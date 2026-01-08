# Wodge UI Overhaul & Asset Management
**Date:** 2026-01-08

## Overview
Major updates focused on modernizing the frontend architecture, implementing a shadcn-inspired component system, and improving the developer onboarding experience with a polished welcome page.

## Key Features

### 1. Wodge UI v2
Moved away from auto-scaffolding all components to an on-demand system. This keeps new projects clean and allows developers to opt-in to complexity.

- **On-Demand Components:** New `wodge add ui <component>` command.
- **Available Components:**
  - `button`: Framer Motion enhanced, highly customizable
  - `card`: Composable layout components
  - `input`: Accessible form controls
  - `navbar`: Responsive navigation
  - `theme-provider`: Dark/light mode context
- **Tech Stack:**
  - React 18 + Vite
  - Tailwind CSS (with consolidated `index.css`)
  - Framer Motion for animations
  - Material UI Icons
  - Assistant font family

### 2. Asset Embedding System
Solved the challenge of distributing binary assets (logos, favicons) within the Go binary.

- **Implementation:** Moved assets to `internal/templates/assets/` to comply with `go:embed` restrictions.
- **Scaffolding:** `wodge new` now automatically:
  - Embeds `logo.png`, `logo_text.png`, and `logo.ico` into the binary.
  - Copies them to the `public/` directory of new projects.
  - Sets `logo.ico` as `public/favicon.ico`.

### 3. Improved Welcome Page
Redesigned the default route to be a comprehensive, yet minimal starting point.

- **Design:** Clean, text-centric layout (Assistant font).
- **Content:**
  - "WODGE" header with logo.
  - Quick start instructions.
  - Feature showcase (UI commands, Backend services, CRUD).
  - Tech stack overview.
- **Developer Experience:** Copy-paste ready commands for exploring features immediately.

### 4. Technical Fixes
- **SSR Hydration:** Fixed `hydrateRoot` vs `createRoot` mismatch in development.
- **Router Architecture:** Moved `BrowserRouter` out of `App.tsx` to entry files to support SSR/SSG correctly.
- **Build System:** Resolved `go:embed` path issues and variable scoping in `root.go`.

## Next Steps
- Expand component library.
- Add more backend service integrations.
- Improve documentation generation.
