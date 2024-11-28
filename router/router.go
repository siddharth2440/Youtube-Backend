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
		// User Registration
		public.POST("/register", authHandlers.RegisterUserHandler)

		// User Login
		public.POST("/login", authHandlers.Login)

		// User Signout
		public.GET("/signout", authHandlers.Signout)

		// Get All Users
		public.GET("/users", userHandlers.GetallUsers)
	}

	// User Routes -> Protected Routes

	private := router.Group("/api/v1/user")
	// Update Profile <- Own Profile
	{
		// Get User Info
		private.GET("/:userID", userHandlers.GetProfile)

		// User Update our profile
		private.PUT("/:userID", userHandlers.UpdateUserHandler)

		// User can delete our profile
		private.DELETE("/:user_id", userHandlers.DeleteUserProfile)

		// User Subscibe
		private.PATCH("/subscribe/:channelid", userHandlers.SubscribeUser)

		// User UnSubscribe
		// private.POST("/:userid", userHandlers.SubscribeUser)
	}

	return router
}
