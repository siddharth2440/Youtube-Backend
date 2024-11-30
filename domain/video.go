package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoStruct struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	UserID    primitive.ObjectID `json:"user_id" validate:"required"`
	Title     string             `json:"title" validate:"required"`
	Desc      string             `json:"description" validate:"required"`
	ImgURL    string             `json:"imgURL,omitempty"`
	VideoURL  string             `json:"videoURL,omitempty"`
	Views     uint64             `json:"views" default:"0"`
	Tags      []*string          `json:"tags" validate:"required"`
	Likes     []string           `json:"likes,omitempty"`
	Dislikes  []string           `json:"dislikes,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func NewVideo(userid *primitive.ObjectID, Title string, Description string, img string, videourl string, tags []*string) *VideoStruct {
	return &VideoStruct{
		ID:        primitive.NewObjectID(),
		UserID:    *userid,
		Title:     Title,
		Desc:      Description,
		ImgURL:    img,
		VideoURL:  videourl,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
