package db

import (
	"context"
	"log"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RetrieveGroupResult int

const (
	Success      RetrieveGroupResult = 0
	Failure      RetrieveGroupResult = 1
	Unauthorised RetrieveGroupResult = 2
)

func GetAllGroups(ctx context.Context) ([]models.Group, bool) {
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
		return nil, false
	}

	var groups []models.Group

	err = cursor.All(ctx, &groups)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return groups, true
}

func GetGroups(ctx context.Context) ([]models.Group, bool) {
	filter := bson.D{{"type", "public"}}

	cursor, err := findGroups(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	var groups []models.Group

	err = cursor.All(ctx, &groups)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return groups, true
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

func GroupExists(ctx context.Context, name string) bool {
	result := findGroup(ctx, name)
	return result.Err() == nil
}

func findGroups(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return GetDatabase().Collection("Groups").Find(ctx, filter)
}

func findGroup(ctx context.Context, name string) *mongo.SingleResult {
	filter := bson.D{{"name", name}}
	return GetDatabase().Collection("Groups").FindOne(ctx, filter)
}

func AddGroup(ctx context.Context, newGroup *models.Group) bool {
	newGroup.ID = uuid.NewString()
	newGroup.TimeCreated = time.Now().UTC().Unix()

	_, err := GetDatabase().Collection("Groups").InsertOne(ctx, newGroup)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func DeleteGroup(ctx context.Context, name string) bool {
	filter := bson.D{{"name", name}}
	_, err := GetDatabase().Collection("Groups").DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
