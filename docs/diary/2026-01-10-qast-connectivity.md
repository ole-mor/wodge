# 2026-01-10: Qast Integration & Networking Fixes

## Overview
Focus was on stabilizing the `wodge` client's connection to the `qast` backend via the proxy, and improving the developer experience for generated apps.

## Key Changes

### 1. QastDriver Improvements
- **Authentication**: Added support for `Authorization: Bearer <token>` in `QastDriver`. 
  - Defaults to `dev-token-bypass` if no API key is provided, enabling frictionless local development against `qast`'s dev-mode AuthGuard.
- **Error Handling**: 
  - `Ask` and `IngestGraph` methods now parse the upstream JSON response body when a non-200 status code is returned. 
  - This exposes actual error messages from `qast` (e.g., specific vector DB or privacy failures) instead of generic "status 500" errors.
- **Network Optimization**: 
  - Configured `http.Client` with a custom `http.Transport`.
  - Enabled **Keep-Alive** and connection pooling (`MaxIdleConns`, `IdleConnTimeout`) to prevent port exhaustion and reduce latency during high-frequency requests (like health checks).

### 2. CLI & Templates
- **Client Generation (`wodge.ts`)**:
  - Updated to use `import.meta.env.VITE_API_URL` for the `API_BASE` path.
  - Defaults to `http://localhost:8080/api` if the env var is missing.
  - Added `VITE_API_URL` to the default `.env` template.
- **UI Components**:
  - `Button`: Added `destructive` variant to `ButtonProps` to support error states in the generated UI.
  - `Card`: Fixed Type errors by using `HTMLMotionProps<'div'>` for proper Framer Motion integration.
  - `QastTest`: Updated the visual feedback component to use these corrected types.

### 3. Distribution
- Verified the `wodge update` command successfully pulls, builds, and installs the latest version from the GitHub repository, establishing a viable update path for the CLI tool.
