# 2026-01-11: Auth Alignment and UI Stability

## Overview
Today's focus was on aligning the authentication flow between the Frontend (`test-app20`), the Proxy (`wodge`), and the Auth Service (`astauth`), along with resolving UI type conflicts.

## Key Changes

### 1. Authentication Flow Alignment (Wodge & AstAuth)
- **Registration & Login**: Updated `wodge` backend and frontend templates to include `username` and `confirmPassword` for registration, aligning with `astauth`'s required DTO fields.
- **Logout Fix**: Resolved a 400/500 error on the logout endpoint.
    - **Issue**: `astauth` required both `access_token` and `refresh_token`, but the frontend was only sending an empty string for the latter.
    - **Fix**: Updated `AuthProvider.tsx` to persist `refresh_token` in `localStorage` and include it in the `logout` API call.
- **Wodge CLI Templates**: Propagated these changes to the `wodge` CLI templates in `components.go` and `add.go` so that all future `wodge add auth` or `wodge new` projects have valid auth logic and corrected `Button` types.

### 2. UI Component Stability
- **Button Type Conflict**: Fixed a recurring TypeScript error in the generated `Button` component in `test-app20`.
    - **Issue**: Conflict between React's `HTMLButtonAttributes` and `framer-motion`'s `HTMLMotionProps` regarding the `onDrag` handler.
    - **Fix**: Updated `ButtonProps` to extend `HTMLMotionProps<'button'>` directly.

### 3. Service Integration
- **Wodge Rebuild**: Successfully rebuilt and reinstalled the `wodge` CLI globally with the latest template fixes.
- **AstAuth Investigation**: Verified service-side validation logic and confirmed DTO requirements for logout tokens.

## Impact
Users can now register, login, and logout successfully across the entire stack without validation errors or TypeScript build failures in the frontend. All future projects scaffolded with Wodge will benefit from these structural stability fixes.
