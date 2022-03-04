package db

import (
	"context"
	"log"
	"phrasmotica/bore-score-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllGroups(ctx context.Context) []models.Group {
	// TODO: filter to global and public groups only
	cursor, err := GetDatabase().Collection("Groups").Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var groups []models.Group

	err = cursor.All(ctx, &groups)
	if err != nil {
		log.Fatal(err)
	}

	return groups
}

func GetGroup(ctx context.Context, name string) models.Group {
	result := findGroup(ctx, name)
	if err := result.Err(); err != nil {
		log.Fatal(err)
	}

	// TODO: return an error if the group is private

	var group models.Group

	if err := result.Decode(&group); err != nil {
		log.Fatal(err)
	}

	return group
}

func findGroup(ctx context.Context, name string) *mongo.SingleResult {
	filter := bson.D{{"name", bson.D{{"$eq", name}}}}
	return GetDatabase().Collection("Groups").FindOne(ctx, filter)
}
