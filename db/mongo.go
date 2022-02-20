package db

import (
	"context"
	"fmt"
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

func GetAllGames(ctx context.Context) []models.Game {
	cursor, err := GetDatabase().Collection("Games").Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var games []models.Game

	err = cursor.All(ctx, &games)
	if err != nil {
		log.Fatal(err)
	}

	return games
}

func GameExists(ctx context.Context, name string) bool {
	filter := bson.D{{"name", bson.D{{"$eq", name}}}}
	result := GetDatabase().Collection("Games").FindOne(ctx, filter)
	return result.Err() == nil
}

func AddGame(ctx context.Context, newGame *models.Game) error {
	newGame.ID = uuid.NewString()
	newGame.Name = computeName(newGame.DisplayName)
	newGame.TimeCreated = time.Now().UTC().Unix()

	_, err := GetDatabase().Collection("Games").InsertOne(ctx, newGame)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeleteGame(ctx context.Context, name string) (bool, error) {
	filter := bson.D{{"name", bson.D{{"$eq", name}}}}
	_, err := GetDatabase().Collection("Games").DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func GetAllPlayers(ctx context.Context) []models.Player {
	cursor, err := GetDatabase().Collection("Players").Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var players []models.Player

	err = cursor.All(ctx, &players)
	if err != nil {
		log.Fatal(err)
	}

	return players
}

func PlayerExists(ctx context.Context, username string) bool {
	filter := bson.D{{"username", bson.D{{"$eq", username}}}}
	result := GetDatabase().Collection("Players").FindOne(ctx, filter)
	return result.Err() == nil
}

func AddPlayer(ctx context.Context, newPlayer *models.Player) error {
	newPlayer.ID = uuid.NewString()
	newPlayer.TimeCreated = time.Now().UTC().Unix()

	_, err := GetDatabase().Collection("Players").InsertOne(ctx, newPlayer)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeletePlayer(ctx context.Context, username string) (bool, error) {
	filter := bson.D{{"username", bson.D{{"$eq", username}}}}
	_, err := GetDatabase().Collection("Players").DeleteOne(ctx, filter)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func GetAllResults(ctx context.Context) []models.Result {
	cursor, err := GetDatabase().Collection("Results").Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []models.Result

	err = cursor.All(ctx, &results)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

func AddResult(ctx context.Context, newResult *models.Result) error {
	newResult.ID = uuid.NewString()
	_, err := GetDatabase().Collection("Results").InsertOne(ctx, newResult)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeleteResultsWithGameId(ctx context.Context, gameId string) (int64, error) {
	filter := bson.D{{"gameId", bson.D{{"$eq", gameId}}}}
	deleteResult, err := GetDatabase().Collection("Results").DeleteMany(ctx, filter)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return deleteResult.DeletedCount, nil
}

func ScrubResultsWithPlayer(ctx context.Context, username string) (int, error) {
	// TODO: do this method via a single collection.Update(...) operation
	results, err := getResultsWithPlayer(ctx, username)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	scrubbedCount := 0

	for i := range results {
		r := results[i]
		for j := range r.Scores {
			s := r.Scores[j]
			if s.Username == username {
				r.Scores[j] = scrubUsernameFromScore(s)
				break
			}
		}

		_, err := GetDatabase().Collection("Results").ReplaceOne(ctx, bson.D{{"id", r.ID}}, r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Scrubbed player %s from result %s\n", username, r.ID)
		scrubbedCount += 1
	}

	return scrubbedCount, nil
}

func getResultsWithPlayer(ctx context.Context, username string) ([]models.Result, error) {
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

	cursor, err := GetDatabase().Collection("Results").Find(ctx, filter)

	if err != nil {
		log.Println(err)
		return []models.Result{}, err
	}

	var results []models.Result

	err = cursor.All(ctx, &results)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Got %d results with player %s\n", len(results), username)

	return results, nil
}

func scrubUsernameFromScore(s models.PlayerScore) models.PlayerScore {
	return models.PlayerScore{
		Username: "",
		Score:    s.Score,
		IsWinner: s.IsWinner,
	}
}
