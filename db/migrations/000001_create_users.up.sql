CREATE TYPE role AS ENUM ('user', 'admin');

CREATE TABLE users (
    id bigint GENERATED ALWAYS AS IDENTITY,
    email varchar(255) NOT NULL,
    name varchar(32) NOT NULL,
    password_hash varchar(255) NOT NULL,
    user_role role NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_id PRIMARY KEY (id),
    CONSTRAINT user_email UNIQUE (email)
);
