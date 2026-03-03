# Regeneration UI and Collaborative Sync
**Date:** 2026-01-28

## Overview
Enhanced the chat interface to support message regeneration, collaborative state synchronization, and direct knowledge inspection.

## Updates

### 1. Regenerate & Inspect Actions
*   **Hover Menu:** Assistant messages now feature an action menu on hover:
    *   **Regenerate (AiEdit icon):** Triggers a new response from the same prompt context.
    *   **Inspect Knowledge:** Provides direct access to the context inspector for the preceding user message.
*   **User Knowledge Inspection:** Added a hover action to user messages that have associated knowledge metadata, allowing for rapid verification/editing of facts.

### 2. Collaborative Sync (History Polling)
*   **Polling Logic:** Implemented 3-second polling for session history.
*   **Highlighting:** When a message is updated by another user or via regeneration, the UI now highlights the changed message with a fading yellow gradient (`highlight-update` class).
*   **Crash Prevention:** Patched `loadHistory` and `pollHistory` to safely handle `null` message responses from the API, preventing frontend crashes during session initialization.

### 3. UI Refinements
*   **Icons:** Integrated `AiEdit.svg` for regeneration and inspection actions.
*   **Theming:** Updated `ContextInspector` aesthetics (colors and borders) to align with the latest design system variables.

## Integration
These features rely on the updated `wodge` backend driver which now supports `target_message_id` and the `UpdateMessage` atomic operation in the Qast service.
