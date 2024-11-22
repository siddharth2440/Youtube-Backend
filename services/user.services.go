package services

import (
	"context"
	"errors"
	"time"

	"github.com/youtube/domain"
	"github.com/youtube/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	db *mongo.Client
}

func NewUserService(db *mongo.Client) *UserService {
	// Initialize the User Collection from the Database
	return &UserService{
		db: db,
	}
}

// Signup Service
func (nus *UserService) RegisterUserService(user *domain.User) (*domain.User, error) {

	// Store the User into The Database
	newUser := domain.NewUser(user.Name, user.Email, user.Password, user.Image)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Insert The User into our Database
	var insertedUser *domain.User

	resultChan := make(chan *domain.User)
	errorChan := make(chan error)

	go func() {
		// Insert the User into the Database
		_, err := nus.db.Database("youtube_db").Collection("users").InsertOne(context.Background(), newUser)
		if err != nil {
			errorChan <- err
			return
		}

		// Retrieve The inseted user
		err = nus.db.Database("youtube_db").Collection("users").FindOne(ctx, bson.M{"_id": newUser.ID}).Decode(&insertedUser)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- insertedUser

	}()

	// Wait for the result, error, or context timeout
	select {
	case user := <-resultChan:
		return user, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

}

// Login Service
func (nus *UserService) Login(email, password *string) (*domain.User, string, error) {
	var user *domain.User
	if email == nil || password == nil {
		return nil, "", errors.New("invalid credentials")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	err := nus.db.Database("youtube_db").Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	defer cancel()

	if err != nil {
		return nil, "", err
	}
	// Check The Password is Valid or Not Valid
	isPasswordValid := utils.ComparePassword(*password, user.Password)

	if !isPasswordValid {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate A JWT Token for further authentication
	token, err := utils.Generate_JWT_Token(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
