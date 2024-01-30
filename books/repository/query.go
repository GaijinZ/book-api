package repository

const (
	SelectAuthor = "SELECT id FROM author WHERE name = $1"
	InsertAuthor = "INSERT INTO author (name) VALUES ($1)"
	InsertBook   = "INSERT INTO user_book (name, date_published, isbn, page_count, user_id, author_id) VALUES ($1, $2, $3, $4, $5, $6)"
	UpdateBook   = "UPDATE user_book SET name = $1, date_published = $2, isbn = $3, page_count = $4, author_id = $5 WHERE id = $6"
	GetBook      = "SELECT b.name, b.date_published, b.isbn, b.page_count, a.name FROM user_book AS JOIN author AS a ON b.author_id = a.id WHERE b.id = $1"
	GetAllBooks  = "SELECT b.name, b.date_published, b.isbn, b.page_count, a.name FROM user_book AS JOIN author AS a ON b.author_id = a.id ORDER BY a.id"
	DeleteBook   = "DELETE FROM user_book WHERE id = $1"
	IsAssigned   = "SELECT EXISTS (SELECT 1 FROM user_book WHERE id = $1 AND user_id = $2)"
	CheckISBN    = "SELECT EXISTS (SELECT 1 FROM user_book WHERE isbn = $1)"
)
