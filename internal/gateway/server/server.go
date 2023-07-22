package server

// import (
// 	"library/database/postgres"
// 	bookHandler "library/internal/books/handler"
// 	bookRepository "library/internal/books/repository"
// 	userHandler "library/internal/users/handler"
// 	userAuth "library/internal/users/repository"
// 	userRepository "library/internal/users/repository"
// 	middleware "library/utils/middleware"

// 	"github.com/gin-gonic/gin"
// )

// func Run(port string) {
// 	dbPool := postgres.GetPostgresConnectionString()
// 	defer dbPool.Close()

// 	bookRepository := bookRepository.NewBookRepository(dbPool)
// 	handlerBook := bookHandler.NewBookHandler(bookRepository)

// 	userRepository := userRepository.NewUserRepository(dbPool)
// 	authRepository := userAuth.NewAuthRepository(dbPool)
// 	authUser := userHandler.NewUserAuth(authRepository)
// 	handlerUser := userHandler.NewUserHandler(userRepository)

// 	router := gin.Default()

// 	router.POST("/v1/users", handlerUser.AddUser)
// 	router.POST("/v1/users/login", authUser.Login)
// 	router.POST("/v1/users/logout", authUser.Logout)

// 	v1 := router.Group("/v1")
// 	v1.Use(middleware.IsAuthorized())

// 	v1.PUT("/users/:user_id", handlerUser.UpdateUser)
// 	v1.GET("/users/:user_id", handlerUser.GetUser)
// 	v1.GET("/users", handlerUser.GetAllUsers)
// 	v1.DELETE("/users/:user_id/:delete_id", handlerUser.DeleteUser)

// 	v1.POST("/users/:user_id/books", handlerBook.AddBook)
// 	v1.PUT("/users/:user_id/books/:book_id", handlerBook.UpdateBook)
// 	v1.GET("/users/:user_id/books/:book_id", handlerBook.GetBook)
// 	v1.GET("/users/:user_id/books", handlerBook.GetAllBooks)
// 	v1.DELETE("/users/:user_id/books/:book_id", handlerBook.DeleteBook)

// 	router.Run(port)
// }
