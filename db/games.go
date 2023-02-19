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

func GetAllGames(ctx context.Context) ([]models.Game, bool) {
	cursor, err := GetDatabase().Collection("Games").Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	var games []models.Game

	err = cursor.All(ctx, &games)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return games, true
}

func GetGame(ctx context.Context, name string) (*models.Game, bool) {
	result := findGame(ctx, name)
	if err := result.Err(); err != nil {
		log.Println(err)
		return nil, false
	}

	var game models.Game

	if err := result.Decode(&game); err != nil {
		log.Println(err)
		return nil, false
	}

	return &game, true
}

func GameExists(ctx context.Context, name string) bool {
	result := findGame(ctx, name)
	return result.Err() == nil
}

func findGame(ctx context.Context, name string) *mongo.SingleResult {
	filter := bson.D{{"name", name}}
	return GetDatabase().Collection("Games").FindOne(ctx, filter)
}

func AddGame(ctx context.Context, newGame *models.Game) bool {
	newGame.ID = uuid.NewString()
	newGame.TimeCreated = time.Now().UTC().Unix()

	_, err := GetDatabase().Collection("Games").InsertOne(ctx, newGame)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func DeleteGame(ctx context.Context, name string) bool {
	filter := bson.D{{"name", name}}
	_, err := GetDatabase().Collection("Games").DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
