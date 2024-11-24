package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youtube/domain"
	"github.com/youtube/services"
)

// User handler struct
type UserHandlerStruct struct {
	userService *services.UserServiceStruct
}

func NewUserHandler(userservice *services.UserServiceStruct) *UserHandlerStruct {
	return &UserHandlerStruct{
		userService: userservice,
	}
}

// Ush := User Service Handler
func (Ush *UserHandlerStruct) UpdateUserHandler(ctx *gin.Context) {

	var user_update_body *domain.User
	userId := ctx.Param("userID")

	if err := ctx.ShouldBindJSON(&user_update_body); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	user_chan := make(chan *domain.User)
	error_chan := make(chan error)

	go func() {

		user, err := Ush.userService.UpdateUserService(user_update_body, &userId)
		if err != nil {
			error_chan <- err
			return
		}
		user_chan <- user
	}()

	select {
	case user := <-user_chan:
		ctx.JSON(http.StatusOK, gin.H{
			"data": user,
		})

	case err := <-error_chan:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func (Ush *UserHandlerStruct) DeleteUserProfile(ctx *gin.Context) {
	userId := ctx.Param("user_id")

	user_channel := make(chan *domain.User)
	err_channel := make(chan error)
	bool_channel := make(chan bool)
	go func() {
		chk, user, err := Ush.userService.DeleteProfile(&userId)
		if err != nil {
			err_channel <- err
			return
		}
		user_channel <- user
		bool_channel <- chk
	}()

	select {
	case user := <-user_channel:
		success := <-bool_channel
		if success {
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"success": success,
					"message": "User deleted successfully",
					"data":    user,
				})
			return
		} else {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{
					"success": success,
					"message": "User is not deleted",
					"data":    nil,
				},
			)
			return
		}
	case err := <-err_channel:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

}
