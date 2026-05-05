# moneyKeeper - Project Goals & Guidelines

## Overview
moneyKeeper is a personal finance management application. It allows users to track their accounts, categorize expenses, and manage transactions between accounts.

## Current Status
- **Backend**: Implemented in Go using the `chi` router.
  - **Database**: PostgreSQL (Recommended).
  - **Functionality**: RESTful API for managing Accounts, Account Types, Categories, and Transactions.
  - **CORS**: Configured to allow requests from `http://localhost:5173`.
  - **Architecture**: Follows standard Go project layout (cmd/internal).
- **Frontend**: Successfully initialized with React and TypeScript in the `frontend/` directory.
  - **Features**: Modern Emerald & Slate UI, persistent layout with navigation, and API-connected Transactions view.
  - **API Integration**: Connected to `/api/account` and `/api/transaction`.

## Project Goals
1. **Frontend Rewrite**: [IN PROGRESS] Rebuild the user interface using React.
2. **Feature Parity**: Ensure the React frontend supports all functionality provided by the backend API.
3. **Data Integrity**: Maintain accurate tracking of financial data across accounts and categories.
4. **User Experience**: Create a clean, intuitive interface for recording and reviewing transactions.

## Technical Stack
- **Backend**: Go (chi, logrus, pq)
- **Database**: PostgreSQL
- **Frontend**: React (TypeScript, Vite, CSS Modules)

## Agent Instructions
- When working on the backend, adhere to the existing Go conventions and project structure.
- When working on the frontend, use CSS Modules for styling and ensure TypeScript types are maintained for API responses.
- Before testing, ensure a PostgreSQL instance is running (e.g., via Docker) and the `backend/dbconfig.yaml` is correctly configured.
