package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youtube/domain"
	"github.com/youtube/response"
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

// Update the User Details
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

// Delete user Profile
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

// Get User Info
func (Ush *UserHandlerStruct) GetProfile(ctx *gin.Context) {
	userID := ctx.Param("userID")

	userChan := make(chan *response.GetUserInfoResponse, 1)
	errChan := make(chan error, 1)

	go func() {
		user, err := Ush.userService.GetUserInfo(&userID)
		if err != nil {
			errChan <- err
			return
		}
		userChan <- user
	}()

	select {
	case user := <-userChan:
		ctx.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	case err := <-errChan:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

// get All Channels
func (Ush *UserHandlerStruct) GetallUsers(ctx *gin.Context) {
	users := make(chan []*response.GetUserInfoResponse, 1)
	errChan := make(chan error, 1)
	go func() {
		usrs, err := Ush.userService.GetAllUsers()
		if err != nil {
			errChan <- err
			return
		}
		users <- usrs
	}()

	select {
	case _users := <-users:
		close(users)
		ctx.JSON(http.StatusOK, gin.H{
			"data": _users,
		})
	case err := <-errChan:
		close(errChan)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

// Subscirbe to the Channel
func (Ush *UserHandlerStruct) SubscribeUser(ctx *gin.Context) {

	channelId := ctx.Param("channelid")

	var RequestBody struct {
		UserID string `json:"userId"`
	}

	if err := ctx.ShouldBindJSON(&RequestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Parsing me error aa rha hai",
		})
		return
	}

	user_subscription_service_data := make(chan *domain.User, 1)
	err_chan := make(chan error, 1)

	go func() {

		userId := RequestBody.UserID

		usr_details, err := Ush.userService.SubscribeToUser(&userId, &channelId)
		if err != nil {
			err_chan <- err
			return
		}
		user_subscription_service_data <- usr_details
	}()

	select {
	case user := <-user_subscription_service_data:
		close(user_subscription_service_data)
		ctx.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	case err := <-err_chan:
		close(err_chan)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}
}

// Unsubscribe to an Channel
