package authent

import (
	"context"
	"errors"
	"fmt"
	"glimpseguru-api/db"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
)

type Identity struct {
	ApiKey string `json:"api_key" binding:"required"`
}

type User struct {
	ID       string    `bson:"userID"`
	APIKey   string    `bson:"apiKey"`
	Websites []Website `bson:"websites"`
}

type Website struct {
	ID   string `bson:"websiteID" json:"websiteID"`
	Name string `bson:"websiteName" json:"websiteName"`
	URL  string `bson:"websiteURL" json:"websiteURL"`
}

func GetUser(apiKey string) (User, error) {
	var user User

	filter := bson.M{
		"apiKey": apiKey,
	}
	result := db.MongoDb.Collection("users").FindOne(context.Background(), filter)
	if result.Err() != nil {
		slog.Error("no user found for apiKey", slog.String("api key", apiKey), slog.String("error", result.Err().Error()))
		return user, errors.New("no user found for apiKey")
	}
	err := result.Decode(&user)
	if err != nil {
		slog.Error(fmt.Sprintf("error decoding user: %e", err))
		return user, errors.New("error decoding user response")
	}

	return user, nil
}
