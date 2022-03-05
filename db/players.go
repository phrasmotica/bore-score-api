package db

import (
	"context"
	"log"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllPlayers(ctx context.Context) ([]models.Player, bool) {
	cursor, err := GetDatabase().Collection("Players").Find(ctx, bson.D{})
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

func PlayerExists(ctx context.Context, username string) bool {
	filter := bson.D{{"username", username}}
	result := GetDatabase().Collection("Players").FindOne(ctx, filter)
	return result.Err() == nil
}

func AddPlayer(ctx context.Context, newPlayer *models.Player) bool {
	newPlayer.ID = uuid.NewString()
	newPlayer.TimeCreated = time.Now().UTC().Unix()

	_, err := GetDatabase().Collection("Players").InsertOne(ctx, newPlayer)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func DeletePlayer(ctx context.Context, username string) bool {
	filter := bson.D{{"username", username}}
	_, err := GetDatabase().Collection("Players").DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
