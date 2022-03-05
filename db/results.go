package db

import (
	"context"
	"log"
	"phrasmotica/bore-score-api/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllResults(ctx context.Context) ([]models.Result, bool) {
	cursor, err := GetDatabase().Collection("Results").Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	var results []models.Result

	err = cursor.All(ctx, &results)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return results, true
}

func AddResult(ctx context.Context, newResult *models.Result) bool {
	newResult.ID = uuid.NewString()

	if len(newResult.GroupName) <= 0 {
		// results are assigned attached to the global group "all" by default
		log.Printf("Assigning new result %s to group all\n", newResult.ID)
		newResult.GroupName = "all"
	}

	_, err := GetDatabase().Collection("Results").InsertOne(ctx, newResult)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func DeleteResultsWithGame(ctx context.Context, gameName string) (int64, bool) {
	filter := bson.D{{"gameName", bson.D{{"$eq", gameName}}}}
	deleteResult, err := GetDatabase().Collection("Results").DeleteMany(ctx, filter)

	if err != nil {
		log.Println(err)
		return 0, false
	}

	return deleteResult.DeletedCount, true
}

func ScrubResultsWithPlayer(ctx context.Context, username string) (int64, bool) {
	// filters to results where the given player took part
	filter := bson.D{
		{
			"scores", bson.D{
				{
					"$elemMatch", bson.D{
						{
							"username", username,
						},
					},
				},
			},
		},
	}

	// updates by setting the username field of the player's score object to an empty string
	update := bson.D{
		{
			"$set", bson.D{
				{
					"scores.$.username", "",
				},
			},
		},
	}

	result, err := GetDatabase().Collection("Results").UpdateMany(ctx, filter, update)

	if err != nil {
		log.Println(err)
		return 0, false
	}

	return result.ModifiedCount, true
}
