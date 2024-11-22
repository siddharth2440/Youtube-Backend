package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youtube/domain"
	"github.com/youtube/services"
)

// import "go.mongodb.org/mongo-driver/mongo"

type AuthHandler struct {
	services *services.UserService
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
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

	user, err := Rh.services.RegisterUserService(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "User registered successfully",
			"user":    user,
		},
	)
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
