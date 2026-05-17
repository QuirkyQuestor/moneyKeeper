# moneyKeeper Backend

RESTful API service for the moneyKeeper finance platform, built with Go and PostgreSQL.

## Core Services
- **Auth**: JWT-based session management with multi-user isolation.
- **Transactions**: Automated double-entry bookkeeping for internal transfers (linked record management).
- **Accounts**: CRUD operations with dynamic balance calculation.
- **Reports**: Aggregation services for income, expenses, and net worth trends.

## Architecture
- `cmd/`: Application entry point.
- `internal/auth`: Authentication logic.
- `internal/datamodel`: Database entity definitions.
- `internal/sqlhandler`: Business logic and database access layer using standard library `database/sql`.

## API Documentation
The API runs on port **8000**. Key endpoints include:
- `GET /api/account/balances` - Fetch summary of all internal account balances.
- `GET /api/transaction` - Fetch paginated transaction records.
- `POST /api/transaction` - Create transaction (automatically creates linked transfers if applicable).
