package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoStruct struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	UserID    string             `json:"user_id" validate:"required"`
	Title     string             `json:"title" validate:"required"`
	Desc      string             `json:"description" validate:"required"`
	ImgURL    string             `json:"imgURL,omitempty"`
	VideoURL  string             `json:"videoURL,omitempty"`
	Views     uint64             `json:"views" default:"0"`
	Tags      []string           `json:"tags" validate:"required"`
	Likes     []string           `json:"likes,omitempty"`
	Dislikes  []string           `json:"dislikes,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func NewVideo(userid, Title, Description, img, videourl string) *VideoStruct {
	return &VideoStruct{
		ID:        primitive.NewObjectID(),
		UserID:    userid,
		Title:     Title,
		Desc:      Description,
		ImgURL:    img,
		VideoURL:  videourl,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
