GRANT ALL PRIVILEGES ON DATABASE postgres TO admin;
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;

CREATE TABLE IF NOT EXISTS public.users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    refresh_token VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    last_login TIMESTAMP,
);

CREATE INDEX idx_email ON users (email);
