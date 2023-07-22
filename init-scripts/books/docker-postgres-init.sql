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
  user_id INTEGER,
  author_id INTEGER REFERENCES authors (id)
);