CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  login VARCHAR(50) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  failed_login_attempts INTEGER DEFAULT 0,
  blocked BOOL DEFAULT false
);

CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  expiration_time timestamptz NOT NULL,
  user_id INTEGER NOT NULL REFERENCES users(id)
);

CREATE TYPE event_type AS ENUM ('login', 'invalid_password', 'block');

CREATE TABLE auth_audit (
  user_id INTEGER NOT NULL REFERENCES users(id),
  event event_type NOT NULL,
  event_time timestamptz DEFAULT NOW()
);