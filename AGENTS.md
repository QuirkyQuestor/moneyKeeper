# moneyKeeper - Project Goals & Guidelines

## Overview
moneyKeeper is a personal finance management application. It allows users to track their accounts, categorize expenses, and manage transactions between accounts.

## Current Status
- **Backend**: Implemented in Go using the `chi` router.
  - **Database**: PostgreSQL 18.3 with UUIDv7 primary keys.
  - **Functionality**: RESTful API for managing Accounts, Account Types, Categories (hierarchical), and Transactions (paginated).
  - **Authentication**: JWT-based multi-user isolation with HttpOnly cookies.
  - **Reporting**: Aggregation endpoints for expenses, income vs. expenses, and net worth trends.
  - **Logging**: Upgraded to standard library `log/slog` for structured JSON logging.
  - **CORS**: Configured to allow requests from `http://localhost:5173`.
  - **Architecture**: Follows standard Go project layout (cmd/internal).
- **Frontend**: Successfully initialized with React and TypeScript.
  - **Features**: Material Design (Material 3) UI, persistent layout, and JWT-based authentication flow.
  - **Authentication**: Integrated with `AuthContext`, `ProtectedRoute`, and a dedicated Login/Register view.
  - **Modal Readiness**: `modal-root` Portal container initialized.
  - **Categories**: Full CRUD with recursive path support.
  - **Transactions**: Paginated table with account filtering and internal/external account classification.
  - **Reports**: Financial dashboard with dynamic charts using `chart.js` and `react-chartjs-2`.

## Development Workflow
- **Orchestration**: Use the root `Makefile` to manage the environment.
- **Database**: PostgreSQL 18.3 runs in Docker via `docker-compose.yml`.
  - Image: `postgres:18.3-alpine`
  - Start DB: `make db-up`
  - Stop DB: `make db-down`
  - Initialize: The schema in `backend/internal/datamodel/schema.sql` is automatically applied on the first `db-up`.
- **Database UI**: pgAdmin 4 is available at `http://localhost:5050`.
  - Login: `admin@moneykeeper.local` / `admin`
- **Backend**: `make run-backend` (runs Go API on port 8000).
- **Frontend**: `make run-frontend` (runs React dev server on port 5173).
- **Full Dev Stack**: `make dev` (starts DB and both apps).

## Project Goals
1. **Frontend Rewrite**: [COMPLETED] Built the user interface using React with a Material Design aesthetic.
2. **Authentication**: [COMPLETED] Multi-user isolation using JWT and HttpOnly cookies.
3. **Feature Parity**: [COMPLETED] React frontend supports all current backend API functionality.
4. **Data Integrity**: Maintain accurate tracking of financial data across accounts and categories.
5. **User Experience**: Maintain a clean, intuitive, and tactile interface for recording and reviewing transactions.

## Technical Stack
- **Backend**: Go (chi, slog, pq, bcrypt, jwt-v5, google/uuid)
- **Database**: PostgreSQL 18.3 (Multi-user schema with user isolation and UUIDv7)
- **Frontend**: React (TypeScript, Vite, CSS Modules, React Router, AuthContext, Chart.js)

## Agent Instructions
- When working on the backend, adhere to the existing Go conventions and project structure.
- When working on the frontend, use CSS Modules for styling, maintain Material Design elevation and radii standards, and utilize React Portals for modals.
- Before testing, ensure the PostgreSQL 18.3 instance is running (via `make db-up`) and the `backend/dbconfig.yaml` is correctly configured.
- Leverage the `is_external` account flag for filtering target accounts in transaction transfers.
- Utilize the hierarchical category path calculation provided by the database-level recursive CTEs in `backend/internal/sqlhandler/category/categoryHelper.go`.
