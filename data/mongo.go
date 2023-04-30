package data

import (
	"context"
	"phrasmotica/bore-score-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	Database *mongo.Database
}

// GetResultsWithPlayer implements IDatabase
func (*MongoDatabase) GetResultsWithPlayer(ctx context.Context, username string) (bool, []models.Result) {
	panic("unimplemented")
}

// UserExists implements IDatabase
func (*MongoDatabase) UserExists(ctx context.Context, email string) bool {
	panic("unimplemented")
}

// GetUser implements IDatabase
func (*MongoDatabase) GetUser(ctx context.Context, username string) (bool, *models.User) {
	panic("unimplemented")
}

// GetUserByEmail implements IDatabase
func (*MongoDatabase) GetUserByEmail(ctx context.Context, email string) (bool, *models.User) {
	panic("unimplemented")
}

// AddUser implements IDatabase
func (*MongoDatabase) AddUser(ctx context.Context, newUser *models.User) bool {
	panic("unimplemented")
}

func CreateMongoDatabase(uri string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		Error.Fatal(err)
		return nil
	}

	return client.Database("BoreScore")
}

func (d *MongoDatabase) GetSummary(ctx context.Context) (bool, *Summary) {
	gameCount, err := d.Database.Collection("Games").CountDocuments(ctx, bson.D{})
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	groupCount, err := d.Database.Collection("Groups").CountDocuments(ctx, bson.D{})
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	playerCount, err := d.Database.Collection("Players").CountDocuments(ctx, bson.D{})
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	resultCount, err := d.Database.Collection("Results").CountDocuments(ctx, bson.D{})
	if err != nil {
		Error.Println(err)
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
		Error.Println(err)
		return false, nil
	}

	var linkTypes []models.LinkType

	err = cursor.All(ctx, &linkTypes)
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	return true, linkTypes
}

func (d *MongoDatabase) GetAllWinMethods(ctx context.Context) (bool, []models.WinMethod) {
	cursor, err := d.Database.Collection("WinMethods").Find(ctx, bson.D{})
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	var winMethods []models.WinMethod

	err = cursor.All(ctx, &winMethods)
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	return true, winMethods
}
