package services

import (
	"context"
	"errors"
	"time"

	"github.com/youtube/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User Services Struct

type UserServiceStruct struct {
	db *mongo.Client
}

func NewUserService(db *mongo.Client) *UserServiceStruct {
	return &UserServiceStruct{
		db: db,
	}
}

// Update User Service  // Nus :--- NewUserService
func (Nus *UserServiceStruct) UpdateUserService(user *domain.User, userId *string) (*domain.User, error) {

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// User Update Channel && Error Channel
	user_update_channel := make(chan *domain.User)
	error_channel := make(chan error)

	go func() {
		// Find the User that it exists OR not
		var is_user_exists *domain.User
		// fmt.Println(*userId)
		// fmt.Println(user)

		userID_in_ObjectID_form, err := primitive.ObjectIDFromHex(*userId)
		if err != nil {
			error_channel <- err
			return
		}
		err = Nus.db.Database("youtube_db").Collection("users").FindOne(ctx, bson.M{"_id": userID_in_ObjectID_form}).Decode(&is_user_exists)
		if err != nil {
			error_channel <- err
			return
		}

		if user.Email != "" {
			is_user_exists.Email = user.Email
		}
		if user.Image != "" {
			is_user_exists.Image = user.Image
		}
		if user.Name != "" {
			is_user_exists.Name = user.Name
		}
		// Update the User Information
		// upsert := true
		// opt := options.UpdateOptions{
		// 	Upsert: &upsert,
		// }

		_, err = Nus.db.Database("youtube_db").Collection("users").UpdateOne(
			ctx,
			bson.M{"_id": userID_in_ObjectID_form},
			bson.D{
				{
					Key: "$set",
					Value: bson.M{
						"email": is_user_exists.Email,
						"image": is_user_exists.Image,
						"name":  is_user_exists.Name,
					},
				},
			},
			// &opt,
		)

		if err != nil {
			error_channel <- err
			return
		}
		user_update_channel <- is_user_exists
		// _ , err := Nus.db.Database("youtube_db").Collection("users").UpdateOne(ctx, bson.M{"_id":"*userID"},bs)
	}()

	select {

	case user_info := <-user_update_channel:
		return user_info, nil

	case get_error := <-error_channel:
		return nil, get_error

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Delete User
func (Nus *UserServiceStruct) DeleteProfile(userId *string) (bool, *domain.User, error) {
	if userId == nil {
		return false, nil, errors.New("invalid user id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	user_delete_channel := make(chan *domain.User)
	errChan := make(chan error)

	user_id, err := primitive.ObjectIDFromHex(*userId)

	if err != nil {
		errChan <- err
		return false, nil, ctx.Err()
	}

	go func() {
		var is_user_exists *domain.User
		err := Nus.db.Database("youtube_db").Collection("users").FindOne(ctx, bson.M{"_id": user_id}).Decode(&is_user_exists)
		if err != nil {
			errChan <- err
			return
		}
		if is_user_exists == nil {
			errChan <- errors.New("user not found")
			return
		}
		user_delete_channel <- is_user_exists

		_, err = Nus.db.Database("youtube_db").Collection("users").DeleteOne(ctx, bson.M{"_id": user_id})
		if err != nil {
			errChan <- err
			return
		}
	}()

	select {
	case delete_user_info := <-user_delete_channel:
		return true, delete_user_info, nil
	case get_error := <-errChan:
		return false, nil, get_error
	case <-ctx.Done():
		return false, nil, ctx.Err()
	}
}

// get User Info
// Subscribe User
// Unsubscribe User

// -- After Video API we write this

// lika a Video
// dislike a Video
