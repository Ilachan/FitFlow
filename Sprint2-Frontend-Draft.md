# Sprint 2 Frontend Draft

## Scope

This file is a standalone frontend draft for Sprint 2 reporting.
It does not modify `Sprint 2.md` and is intended for later merge.

## Frontend Work Completed (Sprint 1 Gap Closure)

### 1. Authentication UX and Validation

- Added real-time form validation for email/password in login and registration flows.
- Added password show/hide toggle in login and registration forms.
- Added password strength indicator and rule checklist in registration page.
- Added registration success confirmation modal while preserving auto-login behavior.

### 2. Profile Editing Flow Completion

- Added profile field-level validation and inline error messages.
- Added avatar URL input and live avatar preview with fallback image.
- Added phone input formatting and validation (`(XXX) XXX-XXXX`).
- Added date of birth age validation (minimum age 13).
- Added save confirmation modal before profile update.
- Added cancel changes action to restore last loaded profile state.
- Added profile refresh after successful save to keep client state aligned with backend data.

### 3. Frontend Quality and Stability

- Added centralized validation utilities for reuse and consistency.
- Fixed frontend lint/build blockers in affected components.
- Verified project compiles and passes lint after changes.

## Current Deferred Item

- Upgrade-to-manager profile action remains a placeholder because backend flow is not fully available yet.

## Frontend Unit Test Implementation

### Test Framework Setup

- Added Vitest and React Testing Library for frontend unit testing.
- Added test scripts in package scripts:
  - `npm test`
  - `npm run test:watch`
- Configured test runtime in `vite.config.ts` (`jsdom` environment + setup file).
- Added `src/test/setup.ts` to enable `jest-dom` and automatic cleanup after each test.

### Unit Test Files Added

- `src/lib/validation.test.ts`
  - Email validation
  - Password policy checks and strength labels
  - Phone formatter and phone validator
  - Avatar URL validator
  - Minimum age validator

- `src/pages/Login.test.tsx`
  - Empty submit shows field validation errors
  - Password show/hide toggle behavior

- `src/pages/Register.test.tsx`
  - Empty submit shows required field errors
  - Password strength indicator updates when user types

- `src/pages/Profile.test.tsx`
  - Invalid avatar URL blocks save and shows inline validation
  - Cancel changes restores original profile values
  - Save confirmation flow calls profile update API on confirm

### Test Execution Result

- Test files: 4 passed
- Tests: 13 passed

### Commands Used

- Run all unit tests:
  - `npm test`
- Run tests in watch mode:
  - `npm run test:watch`
- Validate code quality and build:
  - `npm run lint`
  - `npm run build`

## Suggested Test Placement Convention

- Place reusable utility tests next to utility files (`src/lib/*.test.ts`).
- Place page/component tests next to page/component files (`src/pages/*.test.tsx`, `src/components/*.test.tsx`).
- Keep a shared test setup under `src/test/setup.ts`.
