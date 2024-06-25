-- postgresql

-- tables
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    hash VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    validated_at TIMESTAMP,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS internal_logins (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL references users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    password_salt VARCHAR(255),
    algorithm INTEGER NOT NULL,
    password_last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    login_attempts INT NOT NULL DEFAULT 0,
    last_login_attempt TIMESTAMP,
    last_login_success TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS login_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL references users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    token_type INTEGER NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS external_login_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type INTEGER NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS external_logins (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL references users(id) ON DELETE CASCADE,
    provider_id INT NOT NULL references external_login_providers(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS external_login_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL references users(id) ON DELETE CASCADE,
    provider_id INT NOT NULL references external_login_providers(id) ON DELETE CASCADE,
    login_ip VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    logged_in_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL references users(id) ON DELETE CASCADE,
    login_ip VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    logged_in_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    logged_out_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- indexes
CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS internal_logins_email_idx ON internal_logins (email);
CREATE UNIQUE INDEX IF NOT EXISTS external_login_providers_name_idx ON external_login_providers (name, type, endpoint);
CREATE UNIQUE INDEX IF NOT EXISTS external_logins_user_id_provider_id_idx ON external_logins (user_id, provider_id);
