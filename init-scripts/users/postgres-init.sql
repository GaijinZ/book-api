CREATE USER tmosto WITH PASSWORD 'tmosto' CREATEDB LOGIN;
CREATE DATABASE usersdb WITH OWNER = tmosto;

\c usersdb

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  firstname VARCHAR(100) NOT NULL,
  lastname VARCHAR(100) NOT NULL,
  email VARCHAR(50) NOT NULL,
  password VARCHAR(50) NOT NULL,
  role VARCHAR(10),
  UNIQUE (email)
);

ALTER TABLE users OWNER TO tmosto;