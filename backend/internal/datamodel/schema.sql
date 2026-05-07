-- Create users table
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT uuidv7(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Account types (System provided)
CREATE TABLE account_type (
    type_id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(40) UNIQUE NOT NULL,
    description VARCHAR(200)
);

-- Accounts linked to users
CREATE TABLE account (
    account_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type_id UUID NOT NULL REFERENCES account_type(type_id),
    name VARCHAR(40) NOT NULL,
    description VARCHAR(200),
    active BOOLEAN DEFAULT TRUE,
    is_external BOOLEAN DEFAULT FALSE,
    UNIQUE(user_id, name)
);

-- Categories linked to users
CREATE TABLE category (
    category_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    parent_id UUID REFERENCES category(category_id),
    name VARCHAR(40) NOT NULL,
    description VARCHAR(200),
    expence BOOLEAN DEFAULT TRUE,
    UNIQUE(user_id, name)
);

-- Transactions linked to users
CREATE TABLE transaction (
    transaction_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    account_from UUID NOT NULL REFERENCES account(account_id),
    account_to UUID NOT NULL REFERENCES account(account_id),
    category_id UUID NOT NULL REFERENCES category(category_id),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    memo VARCHAR(255),
    transfer_transaction_id UUID REFERENCES transaction(transaction_id)
);

-- Default account types will be inserted via backend or manual UUID generation
-- Note: PostgreSQL native UUID type is used. IDs should be generated as UUIDv7 in the application layer.
