# 2026-02-08 Deployment Pipeline Fixes

## Summary
Resolved critical issues in the frontend deployment pipeline causing stale code (localhost references) to persist in production.

## Changes

### 1. Docker Build Caching Fix
- Identified that the `test-app20` Dockerfile was copying local `dist/` and `wodge` binaries instead of building from source.
- Implemented a robust build script (`k8s/scripts/build_images.sh`) that:
  - Compiles `wodge` server for Linux/AMD64 locally.
  - Runs `npm install && npm run build` for the frontend locally.
  - Builds the Docker image with fresh artifacts.
- Verified that `localhost:8082` references are replaced by relative `/api` paths in production.

### 2. Wodge Server Logging
- Added explicit request logging middleware to `wodge` server to trace incoming requests.
- This helped debug ingress/routing issues during the 401 troubleshooting.

### 3. API Base URL
- Updated `test-app20/src/lib/wodge.ts` to use relative path `/api` instead of hardcoded `localhost` or complex env logic.
- This ensures seamless routing through the Ingress controller.
