package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youtube/domain"
	"github.com/youtube/services"
)

// import "go.mongodb.org/mongo-driver/mongo"

type AuthHandler struct {
	services *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		services: service,
	}
}

// Signup Handler
func (Rh *AuthHandler) RegisterUserHandler(ctx *gin.Context) {
	var user *domain.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
	}

	// Creating a channels to recieve results from GoRoutine
	resultChan := make(chan *domain.User)
	errChan := make(chan error)

	go func() {
		user, err := Rh.services.RegisterUserService(user)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- user

	}()

	// Wait for result or error
	select {
	case user := <-resultChan:
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": "User registered successfully",
				"user":    user,
			},
		)
	case err := <-errChan:
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": "User registered successfully",
				"error":   err.Error(),
			},
		)
	case <-ctx.Done():
		ctx.JSON(
			http.StatusGatewayTimeout,
			gin.H{
				"error": "Request timed out",
			},
		)
	}
}

// Login Handler
func (Rh *AuthHandler) Login(ctx *gin.Context) {
	var user *domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	user, token, err := Rh.services.Login(&user.Email, &user.Password)
	if err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Invalid email or password" + err.Error(),
			},
		)
		return
	}
	ctx.SetCookie("authCookie", token, 3600, "/", "localhost", false, true)
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Logged in successfully",
			"token":   user,
		},
	)
}

func (Rh *AuthHandler) Signout(ctx *gin.Context) {
	ctx.SetCookie("authCookie", "", -1, "/", "localhost", false, true)
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Logged out successfully",
		},
	)
}
