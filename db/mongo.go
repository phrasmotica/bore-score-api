package db

import (
	"context"
	"log"
	"os"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/google/uuid"
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
		PlayerCount: playerCount,
		ResultCount: resultCount,
	}, true
}

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
	filter := bson.D{{"name", bson.D{{"$eq", name}}}}
	return GetDatabase().Collection("Games").FindOne(ctx, filter)
}

func AddGame(ctx context.Context, newGame *models.Game) bool {
	newGame.ID = uuid.NewString()
	newGame.Name = computeName(newGame.DisplayName)
	newGame.TimeCreated = time.Now().UTC().Unix()

	_, err := GetDatabase().Collection("Games").InsertOne(ctx, newGame)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func DeleteGame(ctx context.Context, name string) bool {
	filter := bson.D{{"name", bson.D{{"$eq", name}}}}
	_, err := GetDatabase().Collection("Games").DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

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
	filter := bson.D{{"username", bson.D{{"$eq", username}}}}
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
	filter := bson.D{{"username", bson.D{{"$eq", username}}}}
	_, err := GetDatabase().Collection("Players").DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

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
