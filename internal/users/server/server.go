package server

import (
	"userapi/internal/users/db/postgres"
	"userapi/internal/users/handler"

	"github.com/gin-gonic/gin"
)

func Run(port string) {
	dbPool := postgres.GetPostgresConnectionString()
	defer dbPool.Close()

	h := handler.DBPoolHandler{DBPool: dbPool}

	router := gin.Default()
	v1 := router.Group("/v1")
	v1.POST("/users", h.AddUser)
	v1.PUT("/users/:id", h.UpdateUser)
	v1.GET("/users/:id", h.GetUser)
	v1.GET("/users", h.GetAllUsers)
	v1.DELETE("/users/:id", h.DeleteUser)

	router.Run(port)
}
