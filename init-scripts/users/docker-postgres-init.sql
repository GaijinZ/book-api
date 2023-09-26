DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'usersdb') THEN
        CREATE DATABASE usersdb WITH OWNER = tmosto;
    END IF;
END $$;

\c booksdb

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  firstname VARCHAR(100) NOT NULL,
  lastname VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL,
  password VARCHAR(200) NOT NULL,
  role VARCHAR(20),
  UNIQUE (email)
);

ALTER TABLE users OWNER TO tmosto;
