CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  login VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(128) NOT NULL,
  failed_login_attempts INTEGER DEFAULT 0,
  blocked BOOL DEFAULT false
);

CREATE TABLE sessions (
  token VARCHAR(32) PRIMARY KEY UNIQUE,
  expiration_time timestamptz NOT NULL
);

CREATE TYPE event_type AS ENUM ('login', 'invalid_password', 'block');

CREATE TABLE auth_audit (
  user_id INTEGER REFERENCES users(id),
  event event_type NOT NULL,
  event_time timestamptz DEFAULT NOW()
);