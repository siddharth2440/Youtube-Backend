package domain

import (
	"time"

	"github.com/youtube/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id" json:"_id"`
	Name            string             `json:"name" validate:"required"`
	Email           string             `json:"email" validate:"required"`
	Password        string             `json:"password" validate:"required"`
	Image           string             `json:"image,omitempty"`
	Subscribers     uint               `json:"subscribers" default:"0"`
	SubscribedUsers []string           `json:"subscribedUsers,omitempty"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

func NewUser(name, email, password, image string) *User {

	hashPassword, _ := utils.HashPassword(password)
	return &User{
		ID:              primitive.NewObjectID(),
		Name:            name,
		Email:           email,
		Password:        hashPassword,
		Image:           image,
		Subscribers:     0,
		SubscribedUsers: []string{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
