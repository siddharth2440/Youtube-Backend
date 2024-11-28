package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youtube/domain"
	"github.com/youtube/services"
)

type VideoHandler struct {
	VideoService *services.VideoService
}

func NewVideoHandler(service *services.VideoService) *VideoHandler {
	return &VideoHandler{VideoService: service}
}

// Define the APIS

// Create a New Video
func (Nvh *VideoHandler) PostAVideo(ctx *gin.Context) {
	var video *domain.VideoStruct
	if err := ctx.ShouldBindJSON(&video); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusCreated,
		gin.H{"message": "Video created successfully", "data": video},
	)
}

// Get a the Video Details
func (Nvh *VideoHandler) GetAVideoInfo(ctx *gin.Context) {
	videoID := ctx.Param("videoID")

	// TODO:
	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "Video details retrieved successfully", "data": videoID},
	)
}

// Update the Video Details
func (Nvh *VideoHandler) UpdateVideoDetails(ctx *gin.Context) {
	var video *domain.VideoStruct
	if err := ctx.ShouldBindJSON(&video); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "Video details updated successfully", "data": video},
	)
}

// Delete the Video Details
func (Nvh *VideoHandler) DeleteTheVideo(ctx *gin.Context) {
	videoID := ctx.Param("videoID")

	// TODO:
	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "Video deleted successfully", "data": videoID},
	)
}

// Get The Random Videos
func (Nvh *VideoHandler) RandomVideos(ctx *gin.Context) {
	//TODO
}

// Add View
func (Nvh *VideoHandler) AddView(ctx *gin.Context) {
	//TODO
}

// Get By Tags
func (Nvh *VideoHandler) GetVideosByTags(ctx *gin.Context) {
	//TODO
}

// Search a Video
func (Nvh *VideoHandler) SearchAVideo(ctx *gin.Context) {
	//TODO
}

// Get All Videos of the Subscribed Channels
func (Nvh *VideoHandler) GetAllVideosOfSubscribedChannels(ctx *gin.Context) {
	//TODO
}
