# 2026-01-06: Wodge CLI Improvements & Refactoring

## Summary
Primary focus was on refining the CLI developer experience (DX), specifically around API generation patterns, command usability (`wodge run`), and observability (`wodge monitor`).

## Key Changes

### 1. CRUD API Refactor (Service Object Pattern)
**Context:** The initial `wodge add api crud` implementation generated "route handlers" (`export func GET`) that returned `Response` objects. This was confusing for a client-side library and didn't fit the Vite + React model well.

**Change:** 
Switched to generating a **Service Object** pattern.
```typescript
import { postgres } from '@/lib/postgres';

export const ItemService = {
  async list() {
    return postgres.query('SELECT * FROM items');
  },
  // ... get, create, delete
};
```
**Impact:** Client code is now cleaner and fully type-safe, directly returning data instead of Response objects.

### 2. Health API
**Context:** Needed a simple, standard way to verify backend connectivity without setting up a database.

**Change:**
Implemented `wodge add api health`.
- Generates `src/api/health.ts`
- Hits the built-in `/api/health` endpoint.
- serves as a minimal example of the `apiGet` helper usage.

### 3. `wodge run` Intelligence
**Context:** `wodge run dev` failed when executed *inside* an app directory because it treated `dev` as the app name.

**Change:**
Updated `wodge run` (in `internal/cli/run.go`) to be context-aware:
- If first arg is a directory -> Switch to it, Run `dev`.
- If first arg is a command (and not a dir) -> Run command in *current* directory.
- If no args -> Run `dev` in *current* directory.

### 4. Monitor CLI Fixes
**Context:** `wodge monitor` was connecting but not displaying logs.

**Fixes:**
- **Event Channel Bug:** Fixed a nil channel initialization issue in the Bubble Tea model that prevented the event stream from actually processing incoming SSE events.
- **Log Formatting:** Updated the TUI table output to mimic Gin's server logs (e.g., `| 200 | GET /path | 12ms |`), making it familiar to Go developers.

## Next Steps
- Implement automated `.env` and `.gitignore` generation for new projects (Planned for Jan 7).
