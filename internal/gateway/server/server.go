package gateway

import (
	"context"
	"library/grpc/proto/user"
	"library/internal/users/handler"
	"log"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// type GRPCHandler struct {
// 	userHandler *handler.UserHandler
// }

// func NewGrpcHandler(userHandler *handler.UserHandler) *GRPCHandler {
// 	return &GRPCHandler{
// 		userHandler: userHandler,
// 	}
// }

// type UserServiceServer struct {
// 	client      user.UserServiceServer
// 	userHandler handler.UserHandler
// }

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	userHandler handler.UserHandler
}

// func GRPCConn() user.UserServiceClient {
// 	// Replace "localhost:5000" with the address where your gRPC server is running.
// 	// You should use the correct IP address and port of your gRPC server.
// 	grpcServerAddress := "localhost:50051"

// 	// Create a gRPC connection to the server
// 	grpcClientConn, err := grpc.Dial(grpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("Failed to connect to the gRPC server: %v", err)
// 	}
// 	defer grpcClientConn.Close()

// 	// Create a gRPC client using the connection
// 	grpcClient := user.NewUserServiceClient(grpcClientConn)

// 	return grpcClient
// }

// Function to handle adding a user via gRPC
func addUserWithGRPC(c *gin.Context) {
	microserviceURL := "http://localhost:5000/v1/users"
	grpcClient := UserServiceServer{}
	var request user.User

	client := http.Client{}

	response, err := client.Get(microserviceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call microservice"})
		return
	}
	defer response.Body.Close()

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	_, err = grpcClient.AddUser(context.Background(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user via gRPC"})
		return
	}
}

func (c *UserServiceServer) AddUser(ctx context.Context, in *user.User) (*user.UserResponse, error) {
	var userResponse user.UserResponse

	firstname := in.FirstName
	lastname := in.LastName
	email := in.Email
	password := in.Password
	role := in.Role

	transformedData := map[string]interface{}{
		"first_name": firstname,
		"last_name":  lastname,
		"email":      email,
		"password":   password,
		"role":       role,
	}

	ginContext := c.createGinContext(transformedData)
	c.userHandler.AddUser(ginContext)

	return &userResponse, nil
}

func (c *UserServiceServer) createGinContext(data map[string]interface{}) *gin.Context {
	// Create a mock HTTP request and context to use within the gin.Context
	httpReq, err := http.NewRequest("POST", "http://localhost:8080", nil)
	if err != nil {
		log.Println("new request error", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Create a *gin.Context object with the mock HTTP request
	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginContext.Request = httpReq

	// Bind the data to the ginContext to be accessed by userHandler.AddUser
	ginContext.Set("transformedData", data)

	return ginContext
}

// Initialize the Gin router and add the route to handle adding a user
func Run() {
	router := gin.Default()

	// Create a gRPC server
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, &UserServiceServer{})

	go func() {
		listenAddress := ":50051"
		lis, err := net.Listen("tcp", listenAddress)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		log.Printf("gRPC server listening on %s", listenAddress)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	microserviceURL := "http://localhost:5000"
	client := http.Client{}
	response, err := client.Get(microserviceURL + "/health")
	if err != nil {
		log.Println("Failed to call microservice")
		return
	}
	defer response.Body.Close()

	// Route for adding a user via gRPC
	router.POST("/v1/users", addUserWithGRPC)

	// Run the Gin server
	router.Run(":8080")
}
