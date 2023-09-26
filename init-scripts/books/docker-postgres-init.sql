DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'booksdb') THEN
        CREATE DATABASE booksdb WITH OWNER = tmosto;
    END IF;
END $$;

\c booksdb

CREATE TABLE IF NOT EXISTS authors (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS book (
  id SERIAL PRIMARY KEY,
  book_name VARCHAR(100) NOT NULL,
  date_published DATE,
  isbn VARCHAR(20),
  page_count INTEGER,
  user_id INTEGER,
  author_id INTEGER REFERENCES authors (id)
);

ALTER TABLE author OWNER TO tmosto;
ALTER TABLE book OWNER TO tmosto;
