CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  firstname VARCHAR(100) NOT NULL,
  lastname VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL,
  password VARCHAR(200) NOT NULL,
  role VARCHAR(20),
  UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS authors (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS books (
  id SERIAL PRIMARY KEY,
  book_name VARCHAR(100) NOT NULL,
  date_published DATE,
  isbn VARCHAR(20),
  page_count INTEGER,
  user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
  author_id INTEGER REFERENCES authors (id)
);