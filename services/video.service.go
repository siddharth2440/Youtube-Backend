package services

import (
	"context"
	"fmt"
	"time"

	"github.com/youtube/domain"
	"github.com/youtube/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoService struct {
	Db       *mongo.Client
	response *response.GetUserInfoResponse
}

func NewVideoService(db *mongo.Client) *VideoService {
	return &VideoService{
		Db:       db,
		response: &response.GetUserInfoResponse{},
	}
}

func (Vs *VideoService) PostAVideo(video *domain.VideoStruct, userid *primitive.ObjectID) (*domain.VideoStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	video_info_chan := make(chan *domain.VideoStruct, 128)
	err_chan := make(chan error, 128)

	created_video := domain.NewVideo(userid, video.Title, video.Desc, video.ImgURL, video.VideoURL, video.Tags)
	go func() {
		info, err := Vs.Db.Database("youtube_db").Collection("videos").InsertOne(ctx, created_video)
		fmt.Println(info)
		if err != nil {
			err_chan <- err
			return
		}
		video_info_chan <- created_video
	}()

	select {
	case video_info := <-video_info_chan:
		close(video_info_chan)
		return video_info, nil
	case err := <-err_chan:
		close(err_chan)
		return nil, err
	}
}

func (Vs *VideoService) GetVideo(videoId *string) (*domain.VideoStruct, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	var videoInfo *domain.VideoStruct
	video := make(chan *domain.VideoStruct, 128)
	err_chan := make(chan error, 128)

	vdio, _ := primitive.ObjectIDFromHex(*videoId)
	fmt.Print(vdio)

	go func() {
		err := Vs.Db.Database("youtube_db").Collection("videos").FindOne(ctx, bson.M{
			"_id": vdio,
		}).Decode(&videoInfo)

		if err != nil {
			err_chan <- err
			return
		}
		video <- videoInfo
	}()

	select {
	case video_info := <-video:
		close(video)
		return video_info, nil
	case err := <-err_chan:
		close(err_chan)
		return nil, err
	}

}

func (Vs *VideoService) UpdateVideo(videoId *string, userId *string, video *domain.VideoStruct) (*domain.VideoStruct, error) {
	var update_the_video *domain.VideoStruct

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	videoId_obj, _ := primitive.ObjectIDFromHex(*videoId)
	userid_obj, _ := primitive.ObjectIDFromHex(*userId)

	update_video_into_channel := make(chan *domain.VideoStruct, 128)
	err_channel := make(chan error, 128)

	filter := bson.D{
		{
			Key: "$and",
			Value: bson.A{
				bson.D{{Key: "userid", Value: userid_obj}},
				bson.D{{Key: "_id", Value: videoId_obj}},
			},
		},
	}
	fmt.Print(videoId_obj)
	fmt.Print(userid_obj)

	go func() {

		// Chk That Video belongs to Our Owner or Not
		err := Vs.Db.Database("youtube_db").Collection("videos").FindOne(ctx, filter).Decode(&update_the_video)
		if err != nil {
			err_channel <- err
			return
		}
		update_video_into_channel <- update_the_video
	}()

	select {
	case video := <-update_video_into_channel:
		close(update_video_into_channel)
		return video, nil
	case err := <-err_channel:
		close(err_channel)
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
