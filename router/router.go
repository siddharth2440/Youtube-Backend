package router

import (
	"github.com/gin-gonic/gin"
	"github.com/youtube/handlers"
	"github.com/youtube/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Client) *gin.Engine {

	router := gin.Default()

	// create Services
	userServices := services.NewUserService(db)

	// create Handlers
	userHandlers := handlers.NewAuthHandler(userServices)

	//Auth Routes
	public := router.Group("/api/v1/auth")
	{
		public.POST("/register", userHandlers.RegisterUserHandler)
		public.POST("/login", userHandlers.Login)
		public.GET("/signout", userHandlers.Signout)
	}

	// User Routes -> Protected Routes

	// Video Routes
	// Comment Routes
	return router
}
