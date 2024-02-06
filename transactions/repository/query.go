package repository

const (
	TransactionHistory        = "SELECT book_list, quantity FROM transactions WHERE user_id = $1"
	CheckAvailability         = "SELECT EXISTS (SELECT 1 FROM book WHERE id = $1 AND quantity > 0)"
	AvailableQuantity         = "SELECT quantity FROM book WHERE id = $1"
	UpdateTransactionQuantity = "UPDATE transactions SET quantity = quantity + $1 WHERE id = $2"
	UpdateBookQuantity        = "UPDATE book SET quantity = quantity - $1 WHERE id = $2"
	GetBook                   = `
								   SELECT
									   book.id,
									   book.name,
									   book.date_published,
									   book.isbn,
									   book.page_count,
									   author.name AS author_name
								   FROM
									   book
									   JOIN book_authors ON book.id = book_authors.book_id
									   JOIN author ON author.id = book_authors.author_id
								   WHERE
									   book.id = $1;
								`
	GetTransactionData = "SELECT id, book_list FROM transactions WHERE user_id = $1"
	InsertTransaction  = "INSERT INTO transactions (user_id, book_list, quantity, transaction_date) VALUES ($1, $2, $3, $4)"
)
