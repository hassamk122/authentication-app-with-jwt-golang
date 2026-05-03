CREATE TABLE IF NOT EXISTS users(
    users_id SERIAL CONSTRAINT users_key PRIMARY KEY,
    public_id UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    username varchar(50) NOT NULL,
    password varchar(255) NOT NULL,
    email varchar(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP 
);