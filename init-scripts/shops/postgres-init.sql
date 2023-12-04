CREATE USER tmosto WITH PASSWORD 'tmosto' CREATEDB LOGIN;
CREATE DATABASE shopsdb WITH OWNER = tmosto;

\c shopsdb

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
    quantity INTEGER,
    author_id INTEGER REFERENCES author (id)
);

CREATE TABLE book_authors (
    book_id INTEGER,
    author_id INTEGER,
    PRIMARY KEY (book_id, author_id),
    FOREIGN KEY (book_id) REFERENCES book(id),
    FOREIGN KEY (author_id) REFERENCES author(id)
);


ALTER TABLE author OWNER TO tmosto;
ALTER TABLE book OWNER TO tmosto;
ALTER TABLE book_authors OWNER TO tmosto;

