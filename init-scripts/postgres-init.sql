CREATE USER tmosto WITH PASSWORD 'tmosto' CREATEDB LOGIN;
CREATE DATABASE booksdb WITH OWNER = tmosto;

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

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id),
    book_list JSON,
    quantity INTEGER,
    transaction_date DATE
);

CREATE TABLE IF NOT EXISTS author (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_book (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    date_published DATE,
    isbn VARCHAR(50),
    page_count INTEGER,
    user_id INTEGER REFERENCES users (id),
    author_id INTEGER REFERENCES author (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS book (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    date_published DATE,
    isbn VARCHAR(50),
    page_count INTEGER,
    quantity INTEGER
);

CREATE TABLE book_authors  (
    book_id INTEGER,
    author_id INTEGER,
    PRIMARY KEY (book_id, author_id),
    FOREIGN KEY (book_id) REFERENCES book(id),
    FOREIGN KEY (author_id) REFERENCES author(id)
);

ALTER TABLE users OWNER TO tmosto;
ALTER TABLE user_book OWNER TO tmosto;
ALTER TABLE transactions OWNER TO tmosto;
ALTER TABLE author OWNER TO tmosto;
ALTER TABLE book OWNER TO tmosto;
ALTER TABLE book_authors  OWNER TO tmosto;
