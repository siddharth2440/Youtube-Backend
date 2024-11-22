package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDB(cfg *Config) (*mongo.Client, error) {
	mongoDBURI := cfg.Mongo_URI

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	optns := options.Client().ApplyURI(mongoDBURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), optns)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	mongoClient := client
	return mongoClient, nil
}

func GetCollections(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("youtube_db").Collection(collectionName)
	return collection
}
