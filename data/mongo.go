package data

import (
	"context"
	"log"
	"os"
	"phrasmotica/bore-score-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	Database *mongo.Database
}

func CreateMongoDatabase() *mongo.Database {
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

func (d *MongoDatabase) GetSummary(ctx context.Context) (bool, *Summary) {
	gameCount, err := d.Database.Collection("Games").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return false, nil
	}

	groupCount, err := d.Database.Collection("Groups").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return false, nil
	}

	playerCount, err := d.Database.Collection("Players").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return false, nil
	}

	resultCount, err := d.Database.Collection("Results").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, &Summary{
		GameCount:   gameCount,
		GroupCount:  groupCount,
		PlayerCount: playerCount,
		ResultCount: resultCount,
	}
}

func (d *MongoDatabase) GetAllLinkTypes(ctx context.Context) (bool, []models.LinkType) {
	cursor, err := d.Database.Collection("LinkTypes").Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return false, nil
	}

	var linkTypes []models.LinkType

	err = cursor.All(ctx, &linkTypes)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, linkTypes
}

func (d *MongoDatabase) GetAllWinMethods(ctx context.Context) (bool, []models.WinMethod) {
	cursor, err := d.Database.Collection("WinMethods").Find(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		return false, nil
	}

	var winMethods []models.WinMethod

	err = cursor.All(ctx, &winMethods)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, winMethods
}
