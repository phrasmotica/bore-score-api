package data

import (
	"context"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *MongoDatabase) GetAllGames(ctx context.Context) (bool, []models.Game) {
	cursor, err := d.Database.Collection("Games").Find(ctx, bson.D{})
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	var games []models.Game

	err = cursor.All(ctx, &games)
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	return true, games
}

func (d *MongoDatabase) GetGame(ctx context.Context, id string) (bool, *models.Game) {
	result := d.findGame(ctx, id)
	if err := result.Err(); err != nil {
		Error.Println(err)
		return false, nil
	}

	var game models.Game

	if err := result.Decode(&game); err != nil {
		Error.Println(err)
		return false, nil
	}

	return true, &game
}

func (d *MongoDatabase) GameExists(ctx context.Context, id string) bool {
	result := d.findGame(ctx, id)
	return result.Err() == nil
}

func (d *MongoDatabase) findGame(ctx context.Context, id string) *mongo.SingleResult {
	filter := bson.D{{"id", id}}
	return d.Database.Collection("Games").FindOne(ctx, filter)
}

func (d *MongoDatabase) AddGame(ctx context.Context, newGame *models.Game) bool {
	newGame.ID = uuid.NewString()
	newGame.TimeCreated = time.Now().UTC().Unix()

	_, err := d.Database.Collection("Games").InsertOne(ctx, newGame)

	if err != nil {
		Error.Println(err)
		return false
	}

	return true
}

func (d *MongoDatabase) DeleteGame(ctx context.Context, id string) bool {
	filter := bson.D{{"id", id}}
	_, err := d.Database.Collection("Games").DeleteOne(ctx, filter)

	if err != nil {
		Error.Println(err)
		return false
	}

	return true
}
