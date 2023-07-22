CREATE USER tmosto WITH PASSWORD 'tmosto' CREATEDB LOGIN;
CREATE DATABASE booksdb WITH OWNER = tmosto;

\c booksdb

CREATE TABLE IF NOT EXISTS author (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS book (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  date_published DATE,
  isbn VARCHAR(50),
  page_count INTEGER,
  user_id INTEGER,
  author_id INTEGER REFERENCES authors (id)
);

ALTER TABLE author OWNER TO tmosto;
ALTER TABLE book OWNER TO tmosto;