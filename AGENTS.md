# moneyKeeper - Project Goals & Guidelines

## Overview
moneyKeeper is a personal finance management application. It allows users to track their accounts, categorize expenses, and manage transactions between accounts.

## Current Status
- **Backend**: Implemented in Go using the `chi` router.
  - **Database**: PostgreSQL with UUIDv7 primary keys.
  - **Functionality**: RESTful API for managing Accounts, Account Types, Categories, and Transactions.
  - **Authentication**: JWT-based multi-user isolation with HttpOnly cookies.
  - **Logging**: Upgraded to standard library `log/slog` for structured JSON logging.
  - **CORS**: Configured to allow requests from `http://localhost:5173`.
  - **Architecture**: Follows standard Go project layout (cmd/internal).
- **Frontend**: Successfully initialized with React and TypeScript.
  - **Features**: Material Design (Material 3) UI, persistent layout, and JWT-based authentication flow.
  - **Authentication**: Integrated with `AuthContext`, `ProtectedRoute`, and a dedicated Login/Register view.
  - **Modal Readiness**: `modal-root` Portal container initialized.
  - **API Integration**: Connected to authenticated endpoints for Accounts and Transactions.

## Development Workflow
- **Orchestration**: Use the root `Makefile` to manage the environment.
- **Database**: PostgreSQL runs in Docker via `docker-compose.yml`.
  - Start DB: `make db-up`
  - Stop DB: `make db-down`
  - Initialize: The schema in `backend/internal/datamodel/schema.sql` is automatically applied on the first `db-up`.
- **Backend**: `make run-backend` (runs Go API on port 8000).
- **Frontend**: `make run-frontend` (runs React dev server on port 5173).
- **Full Dev Stack**: `make dev` (starts DB and both apps).

## Project Goals
1. **Frontend Rewrite**: [IN PROGRESS] Rebuild the user interface using React with a Material Design aesthetic.
2. **Authentication**: [COMPLETED] Multi-user isolation using JWT and HttpOnly cookies.
3. **Feature Parity**: Ensure the React frontend supports all functionality provided by the backend API.
4. **Data Integrity**: Maintain accurate tracking of financial data across accounts and categories.
5. **User Experience**: Create a clean, intuitive, and tactile interface for recording and reviewing transactions.

## Technical Stack
- **Backend**: Go (chi, slog, pq, bcrypt, jwt-v5, google/uuid)
- **Database**: PostgreSQL (Multi-user schema with user isolation and UUIDv7)
- **Frontend**: React (TypeScript, Vite, CSS Modules, React Router, AuthContext)

## Agent Instructions
- When working on the backend, adhere to the existing Go conventions and project structure.
- When working on the frontend, use CSS Modules for styling, maintain Material Design elevation and radii standards, and utilize React Portals for modals.
- Before testing, ensure a PostgreSQL instance is running (e.g., via Docker) and the `backend/dbconfig.yaml` is correctly configured.
