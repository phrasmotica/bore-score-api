package db

import (
	"context"
	"log"
	"phrasmotica/bore-score-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RetrieveGroupResult int

const (
	Success      RetrieveGroupResult = 0
	Failure      RetrieveGroupResult = 1
	Unauthorised RetrieveGroupResult = 2
)

func GetAllGroups(ctx context.Context) *[]models.Group {
	filter := bson.D{
		{
			"type", bson.D{
				{
					"$in", bson.A{"public", "global"},
				},
			},
		},
	}

	cursor, err := findGroups(ctx, filter)
	if err != nil {
		log.Println(err)
		return &[]models.Group{}
	}

	var groups []models.Group

	err = cursor.All(ctx, &groups)
	if err != nil {
		log.Println(err)
		return &[]models.Group{}
	}

	return &groups
}

func GetGroup(ctx context.Context, name string) (*models.Group, RetrieveGroupResult) {
	result := findGroup(ctx, name)
	if err := result.Err(); err != nil {
		log.Println(err)
		return nil, Failure
	}

	var group models.Group

	if err := result.Decode(&group); err != nil {
		log.Println(err)
		return nil, Failure
	}

	if group.Type == models.Private {
		return nil, Unauthorised
	}

	return &group, Success
}

func findGroups(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return GetDatabase().Collection("Groups").Find(ctx, filter)
}

func findGroup(ctx context.Context, name string) *mongo.SingleResult {
	filter := bson.D{{"name", bson.D{{"$eq", name}}}}
	return GetDatabase().Collection("Groups").FindOne(ctx, filter)
}
