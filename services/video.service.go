package services

import "go.mongodb.org/mongo-driver/mongo"

type VideoService struct {
	Db *mongo.Client
	// TODO :- Add Response Type in future
}

func NewVideoService(db *mongo.Client) *VideoService {
	return &VideoService{Db: db}
}
