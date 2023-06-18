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
	v1.POST("/add-user", h.AddUser)
	v1.PUT("/update-user/:id", h.UpdateUser)
	v1.GET("/get-user/:id", h.GetUser)
	v1.GET("/get-users", h.GetAllUsers)
	v1.DELETE("/delete-user/:id", h.DeleteUser)

	router.Run(port)
}
