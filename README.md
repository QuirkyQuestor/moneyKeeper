# moneyKeeper

A modern, tactile personal finance management application designed for tracking accounts, categories, and financial transactions with multi-user isolation.

## Key Features
- **Automated Transfers**: Double-entry ledger system for internal account transfers.
- **Home Dashboard**: Real-time net worth tracking and account balance overview.
- **Hierarchical Categories**: Organize expenses and income with recursive path support.
- **Transaction Ledger**: Paginated, filterable transaction management with real-time balance updates.
- **Secure Auth**: Multi-user isolation powered by JWT and HttpOnly cookies.

## Tech Stack
- **Backend**: Go (chi router, log/slog)
- **Frontend**: React, TypeScript, Vite, CSS Modules, Material 3
- **Database**: PostgreSQL 18.3 (with UUIDv7)

## Getting Started
- **Development Stack**: `make dev` (starts database, backend API, and frontend dev server).
- **Database UI**: pgAdmin is no longer included; use your preferred tool to connect to port 5432.
