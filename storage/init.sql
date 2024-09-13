CREATE TABLE users (
                       user_id UUID PRIMARY KEY,
                       login VARCHAR(50) UNIQUE NOT NULL,
                       password TEXT NOT NULL
);

CREATE TABLE tasks (
                       task_id UUID PRIMARY KEY,
                       user_id UUID NOT NULL REFERENCES users(user_id),
                       payload JSONB NOT NULL,
                       status VARCHAR(50) NOT NULL,
                       result TEXT DEFAULT NULL
);

CREATE INDEX idx_users_login ON users(login);