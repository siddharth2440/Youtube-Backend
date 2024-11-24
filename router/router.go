package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/youtube/handlers"
	"github.com/youtube/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Client) *gin.Engine {

	router := gin.Default()

	// -- CORS -- configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	// create Services
	authServices := services.NewAuthService(db)
	userServices := services.NewUserService(db)

	// create Handlers
	authHandlers := handlers.NewAuthHandler(authServices)
	userHandlers := handlers.NewUserHandler(userServices)

	//Auth Routes
	public := router.Group("/api/v1/auth")
	{
		public.POST("/register", authHandlers.RegisterUserHandler)
		public.POST("/login", authHandlers.Login)
		public.GET("/signout", authHandlers.Signout)
	}

	// User Routes -> Protected Routes

	private := router.Group("/api/v1/user")
	// Update Profile <- Own Profile
	{
		private.PUT("/:userID", userHandlers.UpdateUserHandler)
		private.DELETE("/:user_id", userHandlers.DeleteUserProfile)
	}
	// Delete Profile <- Own Profile
	// get User Info
	// Subscribe User
	// Unsubscribe User
	// like a video
	// dislike a video

	return router
}
