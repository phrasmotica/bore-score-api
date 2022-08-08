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

func players() *mongo.Collection {
	return GetDatabase().Collection("Players")
}

func GetAllPlayers(ctx context.Context) ([]models.Player, bool) {
	cursor, err := players().Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	var players []models.Player

	err = cursor.All(ctx, &players)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return players, true
}

func GetPlayer(ctx context.Context, username string) (*models.Player, bool) {
	result := findPlayer(ctx, username)
	if err := result.Err(); err != nil {
		log.Println(err)
		return nil, false
	}

	var player models.Player

	if err := result.Decode(&player); err != nil {
		log.Println(err)
		return nil, false
	}

	return &player, true
}

func PlayerExists(ctx context.Context, username string) bool {
	result := findPlayer(ctx, username)
	return result.Err() == nil
}

func findPlayer(ctx context.Context, username string) *mongo.SingleResult {
	filter := bson.D{{"username", username}}
	return players().FindOne(ctx, filter)
}

func AddPlayer(ctx context.Context, newPlayer *models.Player) bool {
	newPlayer.ID = uuid.NewString()
	newPlayer.TimeCreated = time.Now().UTC().Unix()

	_, err := players().InsertOne(ctx, newPlayer)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func DeletePlayer(ctx context.Context, username string) bool {
	filter := bson.D{{"username", username}}
	_, err := players().DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
