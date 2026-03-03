# Context Update Response Passthrough
**Date:** 2026-02-01

## Overview
Fixed the wodge gateway to properly pass through the full response from qast's context update endpoint, enabling the frontend to receive entity token mappings.

## The Problem
The Context Inspector's edit feature was receiving `{"status": "updated"}` instead of the full response containing `graph` and `token_map`. This prevented the frontend from resolving entity tokens to display actual names.

## Root Cause
The `handleContextUpdate` handler in the wodge server was calling `qastSvc.UpdateContext()` but:
1. The `UpdateContext` method only returned `error`, discarding qast's response
2. The handler returned a hardcoded `{"status": "updated"}` instead of forwarding the response

## Fixes Applied

### 1. QastDriver.UpdateContext
*   **File:** `internal/drivers/qast/driver.go`
*   Changed return type from `error` to `(map[string]interface{}, error)`
*   Now parses and returns the full JSON response from qast (includes `graph`, `token_map`, `message`, `id`)

### 2. QastService Interface
*   **File:** `internal/services/interfaces.go`
*   Updated `UpdateContext` signature to match the new return type

### 3. handleContextUpdate Handler
*   **File:** `internal/server/server.go`
*   Now captures the response from `UpdateContext` and passes it through to the frontend
*   Removed hardcoded `{"status": "updated"}` response

## Result
Frontend now receives the complete response including `token_map`, allowing `TokenManager` to save entity mappings and resolve tokens to display actual names in the Context Inspector.
