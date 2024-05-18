CREATE TYPE user_type AS ENUM('IT', 'Nurse');

CREATE TABLE IF NOT EXISTS
users (
	id CHAR(16) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    nip INT(13) NOT NULL UNIQUE,
    user_type user_type NOT NULL,
    hashed_password BYTEA,
    identity_card_url TEXT,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS users_nip
	ON users(nip);
CREATE INDEX IF NOT EXISTS users_user_type
	ON users(user_type);
CREATE INDEX IF NOT EXISTS users_name
	ON users(lower(name));
CREATE INDEX IF NOT EXISTS users_created_at_desc
	ON users(created_at DESC);
CREATE INDEX IF NOT EXISTS users_created_at_asc
	ON users(created_at ASC);