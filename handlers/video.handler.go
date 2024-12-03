package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/youtube/domain"
	"github.com/youtube/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	video_info := make(chan *domain.VideoStruct, 128)
	err_chan := make(chan error, 128)

	userid := ctx.GetString("userId")

	startIndex := strings.Index(userid, "\"") + 1
	endIndex := strings.LastIndex(userid, "\"")
	objecid := userid[startIndex:endIndex]
	userID, _ := primitive.ObjectIDFromHex(objecid)

	// fmt.Println("Converted ObjectID", userID)

	go func() {
		video, err := Nvh.VideoService.PostAVideo(video, &userID)
		if err != nil {
			err_chan <- err
			return
		}
		video_info <- video
	}()

	select {
	case video := <-video_info:
		close(video_info)
		ctx.JSON(
			http.StatusCreated,
			gin.H{"message": "Video created successfully", "data": video},
		)
	case err := <-err_chan:
		close(err_chan)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

}

// Get a the Video Details
func (Nvh *VideoHandler) GetAVideoInfo(ctx *gin.Context) {
	videoID := ctx.Param("videoId")
	fmt.Println(videoID)
	video := make(chan *domain.VideoStruct, 128)
	err_chan := make(chan error, 128)

	go func() {
		video_info, err := Nvh.VideoService.GetVideo(&videoID)
		if err != nil {
			err_chan <- err
			return
		}
		video <- video_info
	}()

	select {
	case _video := <-video:
		close(video)
		ctx.JSON(
			http.StatusOK,
			gin.H{"message": "Video details retrieved successfully", "data": _video},
		)
	case err := <-err_chan:
		close(err_chan)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Update the Video Details
func (Nvh *VideoHandler) UpdateVideoDetails(ctx *gin.Context) {
	var video *domain.VideoStruct
	videoID := ctx.Param("videoId")
	fmt.Println("Pehle :- ", videoID)
	if err := ctx.ShouldBindJSON(&video); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("baad me :- ", videoID)

	video_chan := make(chan *domain.VideoStruct, 128)
	err_chan := make(chan error, 128)
	userID := ctx.GetString("userId")
	firstIndex := strings.Index(userID, "\"") + 1
	lastIndex := strings.LastIndex(userID, "\"")

	userID = userID[firstIndex:lastIndex]
	fmt.Println(userID)
	fmt.Println(videoID)

	go func() {
		get_video, err := Nvh.VideoService.UpdateVideo(&videoID, &userID, video)
		if err != nil {
			err_chan <- err
			return
		}
		video_chan <- get_video
	}()

	select {
	case video := <-video_chan:
		close(video_chan)
		ctx.JSON(
			http.StatusOK,
			gin.H{"message": "Video details updated successfully", "data": video},
		)
	case err := <-err_chan:
		close(err_chan)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Unable to update video details", "error": err.Error()},
		)
	}
}

// Delete the Video Details   ---->
func (Nvh *VideoHandler) DeleteTheVideo(ctx *gin.Context) {
	videoID := ctx.Param("videoID")

	// TODO:
	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "Video deleted successfully", "data": videoID},
	)
}

// Get The Random Videos  ---->
func (Nvh *VideoHandler) RandomVideos(ctx *gin.Context) {
	//TODO
}

// Add View
func (Nvh *VideoHandler) AddView(ctx *gin.Context) {
	//TODO
}

// Get By Tags ---->
func (Nvh *VideoHandler) GetVideosByTags(ctx *gin.Context) {
	//TODO
}

// Search a Video ---->
func (Nvh *VideoHandler) SearchAVideo(ctx *gin.Context) {
	//TODO
}

// Get All Videos of the Subscribed Channels
func (Nvh *VideoHandler) GetAllVideosOfSubscribedChannels(ctx *gin.Context) {
	//TODO
}
