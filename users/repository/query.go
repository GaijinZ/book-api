package repository

const (
	GetUserByEmail = "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"
	InsertUser     = `
					INSERT INTO users (firstname, lastname, email, password, role) 
					VALUES ($1, $2, $3, $4, $5) 
					`
	UpdateUser  = "UPDATE users SET firstname=$1, lastname=$2, email=$3, role=$4 WHERE id=$5"
	GetUserByID = "SELECT firstname, lastname, email, role FROM users WHERE id=$1"
	GetUsers    = "SELECT firstname, lastname, email, role FROM users ORDER BY id"
	DeleteUser  = "DELETE FROM users WHERE id=$1"
)
