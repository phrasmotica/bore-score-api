package data

import (
	"context"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *MongoDatabase) GetAllResults(ctx context.Context) (bool, []models.Result) {
	cursor, err := d.Database.Collection("Results").Find(ctx, bson.D{})
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	var results []models.Result

	err = cursor.All(ctx, &results)
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	return true, results
}

func (d *MongoDatabase) AddResult(ctx context.Context, newResult *models.Result) bool {
	newResult.ID = uuid.NewString()
	newResult.TimeCreated = time.Now().UTC().Unix()

	_, err := d.Database.Collection("Results").InsertOne(ctx, newResult)

	if err != nil {
		Error.Println(err)
		return false
	}

	return true
}

func (d *MongoDatabase) DeleteResultsWithGame(ctx context.Context, gameId string) (bool, int64) {
	filter := bson.D{{"gameId", gameId}}
	deleteResult, err := d.Database.Collection("Results").DeleteMany(ctx, filter)

	if err != nil {
		Error.Println(err)
		return false, 0
	}

	return true, deleteResult.DeletedCount
}

func (d *MongoDatabase) ScrubResultsWithPlayer(ctx context.Context, username string) (bool, int64) {
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

	result, err := d.Database.Collection("Results").UpdateMany(ctx, filter, update)

	if err != nil {
		Error.Println(err)
		return false, 0
	}

	return true, result.ModifiedCount
}
