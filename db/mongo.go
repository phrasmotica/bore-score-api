package db

import (
	"context"
	"log"
	"os"
	"phrasmotica/bore-score-api/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

func connect() *mongo.Database {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("No .env.local file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("No MONGODB_URI environment variable found!")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client.Database("BoreScore")
}

func GetDatabase() *mongo.Database {
	if database == nil {
		database = connect()
	}

	return database
}

func GetSummary(ctx context.Context) (*Summary, bool) {
	gameCount, err := GetDatabase().Collection("Games").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	groupCount, err := GetDatabase().Collection("Groups").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	playerCount, err := GetDatabase().Collection("Players").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	resultCount, err := GetDatabase().Collection("Results").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return &Summary{
		GameCount:   gameCount,
		GroupCount:  groupCount,
		PlayerCount: playerCount,
		ResultCount: resultCount,
	}, true
}

func GetAllLinkTypes(ctx context.Context) ([]models.LinkType, bool) {
	cursor, err := GetDatabase().Collection("LinkTypes").Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	var linkTypes []models.LinkType

	err = cursor.All(ctx, &linkTypes)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return linkTypes, true
}

func GetAllWinMethods(ctx context.Context) ([]models.WinMethod, bool) {
	cursor, err := GetDatabase().Collection("WinMethods").Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	var winMethods []models.WinMethod

	err = cursor.All(ctx, &winMethods)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return winMethods, true
}
