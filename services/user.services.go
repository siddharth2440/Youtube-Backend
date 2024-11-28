package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/youtube/domain"
	"github.com/youtube/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User Services Struct

type UserServiceStruct struct {
	db  *mongo.Client
	res *response.GetUserInfoResponse
}

func NewUserService(db *mongo.Client) *UserServiceStruct {
	return &UserServiceStruct{
		db:  db,
		res: &response.GetUserInfoResponse{},
	}
}

// Update User Service  // Nus :--- NewUserService
func (Nus *UserServiceStruct) UpdateUserService(user *domain.User, userId *string) (*domain.User, error) {

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// User Update Channel && Error Channel
	user_update_channel := make(chan *domain.User, 1)
	error_channel := make(chan error, 1)

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
		close(user_update_channel)
		return user_info, nil

	case get_error := <-error_channel:
		close(error_channel)
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

	user_delete_channel := make(chan *domain.User, 1)
	errChan := make(chan error, 1)

	user_id, err := primitive.ObjectIDFromHex(*userId)
	fmt.Print(user_id)

	// deleteeRes := make(chan any, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	if err != nil {
		errChan <- err
		return false, nil, ctx.Err()
	}
	fmt.Println("Deleteting the UserId :", *userId)

	go func() {
		defer wg.Done()
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

	wg.Wait()

	select {
	case delete_user_info := <-user_delete_channel:
		close(user_delete_channel)
		// print(deleteeRes)
		return true, delete_user_info, nil
	case get_error := <-errChan:
		close(errChan)
		return false, nil, get_error
	case <-ctx.Done():
		return false, nil, ctx.Err()
	}
}

// get User Info
func (Nus *UserServiceStruct) GetUserInfo(userId *string) (*response.GetUserInfoResponse, error) {
	var user *domain.User
	if userId == nil {
		return nil, errors.New("userid is required")
	}

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	userid_to_objectID, _ := primitive.ObjectIDFromHex(*userId)

	userInfo := make(chan *response.GetUserInfoResponse, 1)
	errorChan := make(chan error, 1)

	filter := bson.M{
		"_id": userid_to_objectID,
	}

	go func() {
		// Find the User In the Database
		err := Nus.db.Database("youtube_db").Collection("users").FindOne(ctx, filter).Decode(&user)
		if err != nil {
			errorChan <- err
			return
		}
		userDetails := Nus.res.GetDetails(user)
		userInfo <- userDetails
	}()

	select {

	case user := <-userInfo:
		close(userInfo)
		return user, nil
	case err := <-errorChan:
		close(errorChan)
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// get All Users
func (Nus *UserServiceStruct) GetAllUsers() ([]*response.GetUserInfoResponse, error) {

	var users []*domain.User
	allChannels := make(chan []*response.GetUserInfoResponse, 1)
	err_chan := make(chan error, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	go func() {
		cur, err := Nus.db.Database("youtube_db").Collection("users").Find(ctx, bson.M{})
		if err != nil {
			err_chan <- err
			return
		}

		for cur.Next(ctx) {
			var user *domain.User
			err := cur.Decode(&user)
			if err != nil {
				err_chan <- err
				return
			}
			// fmt.Printf("%v\n", user)
			users = append(users, user)
		}
		users_as_response := Nus.res.GetAllUsersResponse(users)
		allChannels <- users_as_response
		defer cur.Close(ctx)
	}()

	select {
	case users := <-allChannels:
		close(allChannels)
		return users, nil
	case err := <-err_chan:
		close(err_chan)
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Subscribe User
func (Nus *UserServiceStruct) SubscribeToUser(userId, channelId *string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	var user *domain.User
	if userId == nil || channelId == nil {
		return nil, errors.New("error in userid mismatching or null")
	}

	// WaitGroups
	var wg sync.WaitGroup
	wg.Add(2)

	user_ObjectID_from_hex, err := primitive.ObjectIDFromHex(*userId)
	if err != nil {
		return nil, err
	}
	channel_ObjectID_from_hex, err := primitive.ObjectIDFromHex(*channelId)
	if err != nil {
		return nil, err
	}

	// We can't subscribe to Ourselves
	filterToUpdateTheUser := bson.M{
		"_id": user_ObjectID_from_hex,
	}
	filterToUpdateTheChannel := bson.M{
		"_id": channel_ObjectID_from_hex,
	}

	// Object
	updateUser := bson.M{
		"$addToSet": bson.M{
			"subscribedusers": channelId,
		},
	}
	updateChannel := bson.M{
		"$inc": bson.M{
			"subscribers": 1,
		},
	}

	fmt.Print("Yahan toh Aaya Service pe")
	get_info_after_update := make(chan *domain.User, 1)
	err_after_update := make(chan error, 1)

	// Go Routine to Update the User
	go func() {
		defer wg.Done()

		// Add the Subscription for the User to the Channel
		err = Nus.db.Database("youtube_db").Collection("users").FindOneAndUpdate(ctx, filterToUpdateTheUser, updateUser).Decode(&user)
		if err != nil {
			err_after_update <- err
			return
		}
		get_info_after_update <- user
	}()

	// Go Routine to Update the Chennel
	go func() {
		defer wg.Done()

		// Increment Channels Subscribers Count
		err := Nus.db.Database("youtube_db").Collection("users").FindOneAndUpdate(ctx, filterToUpdateTheChannel, updateChannel).Err()
		if err != nil {
			err_after_update <- err
			return
		}
	}()

	// Wait for both updates to complete
	wg.Wait()

	select {
	case user := <-get_info_after_update:
		close(get_info_after_update)
		return user, nil
	case err := <-err_after_update:
		close(err_after_update)
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Unsubscribe User
func (Nus *UserServiceStruct) UnSubscribeToUser(userId, channelId *string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	var user *domain.User
	if userId == nil || channelId == nil {
		return nil, errors.New("error in userid mismatching or null")
	}

	if *userId == *channelId {
		return nil, errors.New("error in userid mismatching or null")
	}

	user_ObjectID_from_hex, err := primitive.ObjectIDFromHex(*userId)
	if err != nil {
		return nil, err
	}
	channel_ObjectID_from_hex, err := primitive.ObjectIDFromHex(*channelId)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	filterToUpdateTheUser := bson.M{
		"_id": user_ObjectID_from_hex,
	}
	filterToUpdateTheChannel := bson.M{
		"_id": channel_ObjectID_from_hex,
	}

	// Object
	updateUser := bson.M{
		"$pull": bson.M{
			"subscribedusers": channelId,
		},
	}
	updateChannel := bson.M{
		"$inc": bson.M{
			"subscribers": -1,
		},
	}

	get_info_after_update := make(chan *domain.User, 1)
	err_after_update := make(chan error, 1)

	// Go Routine to Update the User
	go func() {
		defer wg.Done()
		err := Nus.db.Database("youtube_id").Collection("users").FindOne(ctx, bson.M{"_id": user_ObjectID_from_hex}).Decode(&user)

		if err != nil {
			err_after_update <- err
			return
		}

		get_info_after_update <- user

		// Add the Subscription for the User to the Channel
		err = Nus.db.Database("youtube_id").Collection("users").FindOneAndUpdate(ctx, filterToUpdateTheUser, updateUser).Decode(&user)
		if err != nil {
			err_after_update <- err
			return
		}
		get_info_after_update <- user
	}()

	// Go Routine to Update the Chennel
	go func() {
		defer wg.Done()
		// Increment Channels Subscribers Count
		err := Nus.db.Database("youtube_id").Collection("users").FindOneAndUpdate(ctx, filterToUpdateTheChannel, updateChannel).Err()
		if err != nil {
			err_after_update <- err
			return
		}
	}()

	go func() {
		wg.Wait()
		close(get_info_after_update)
		close(err_after_update)
	}()

	select {
	case user := <-get_info_after_update:
		return user, nil
	case err := <-err_after_update:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// -- After Video API we write this

// lika a Video
// dislike a Video
