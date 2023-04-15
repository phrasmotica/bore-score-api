package data

import (
	"context"
	"log"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *MongoDatabase) players() *mongo.Collection {
	return d.Database.Collection("Players")
}

func (d *MongoDatabase) GetAllPlayers(ctx context.Context) (bool, []models.Player) {
	cursor, err := d.players().Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return false, nil
	}

	var players []models.Player

	err = cursor.All(ctx, &players)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, players
}

func (d *MongoDatabase) GetPlayer(ctx context.Context, username string) (bool, *models.Player) {
	result := d.findPlayer(ctx, username)
	if err := result.Err(); err != nil {
		log.Println(err)
		return false, nil
	}

	var player models.Player

	if err := result.Decode(&player); err != nil {
		log.Println(err)
		return false, nil
	}

	return true, &player
}

func (d *MongoDatabase) PlayerExists(ctx context.Context, username string) bool {
	result := d.findPlayer(ctx, username)
	return result.Err() == nil
}

func (d *MongoDatabase) findPlayer(ctx context.Context, username string) *mongo.SingleResult {
	filter := bson.D{{"username", username}}
	return d.players().FindOne(ctx, filter)
}

func (d *MongoDatabase) AddPlayer(ctx context.Context, newPlayer *models.Player) bool {
	newPlayer.ID = uuid.NewString()
	newPlayer.TimeCreated = time.Now().UTC().Unix()

	_, err := d.players().InsertOne(ctx, newPlayer)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (d *MongoDatabase) DeletePlayer(ctx context.Context, username string) bool {
	filter := bson.D{{"username", username}}
	_, err := d.players().DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
