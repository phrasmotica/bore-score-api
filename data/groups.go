package data

import (
	"context"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *MongoDatabase) GetAllGroups(ctx context.Context) (bool, []models.Group) {
	filter := bson.D{
		{
			"visibility", bson.D{
				{
					"$in", bson.A{"public", "global"},
				},
			},
		},
	}

	cursor, err := d.findGroups(ctx, filter)
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	var groups []models.Group

	err = cursor.All(ctx, &groups)
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	return true, groups
}

func (d *MongoDatabase) GetGroups(ctx context.Context) (bool, []models.Group) {
	filter := bson.D{{"visibility", "public"}}

	cursor, err := d.findGroups(ctx, filter)
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	var groups []models.Group

	err = cursor.All(ctx, &groups)
	if err != nil {
		Error.Println(err)
		return false, nil
	}

	return true, groups
}

func (d *MongoDatabase) GetGroup(ctx context.Context, id string) (bool, *models.Group) {
	result := d.findGroup(ctx, id)
	if err := result.Err(); err != nil {
		Error.Println(err)
		return false, nil
	}

	var group models.Group

	if err := result.Decode(&group); err != nil {
		Error.Println(err)
		return false, nil
	}

	return true, &group
}

func (d *MongoDatabase) GetGroupByName(ctx context.Context, name string) (bool, *models.Group) {
	result := d.findGroupByName(ctx, name)
	if err := result.Err(); err != nil {
		Error.Println(err)
		return false, nil
	}

	var group models.Group

	if err := result.Decode(&group); err != nil {
		Error.Println(err)
		return false, nil
	}

	return true, &group
}

func (d *MongoDatabase) GroupExists(ctx context.Context, name string) bool {
	result := d.findGroupByName(ctx, name)
	return result.Err() == nil
}

func (d *MongoDatabase) findGroups(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return d.Database.Collection("Groups").Find(ctx, filter)
}

func (d *MongoDatabase) findGroup(ctx context.Context, id string) *mongo.SingleResult {
	filter := bson.D{{"id", id}}
	return d.Database.Collection("Groups").FindOne(ctx, filter)
}

func (d *MongoDatabase) findGroupByName(ctx context.Context, name string) *mongo.SingleResult {
	filter := bson.D{{"name", name}}
	return d.Database.Collection("Groups").FindOne(ctx, filter)
}

func (d *MongoDatabase) AddGroup(ctx context.Context, newGroup *models.Group) bool {
	newGroup.ID = uuid.NewString()
	newGroup.TimeCreated = time.Now().UTC().Unix()

	_, err := d.Database.Collection("Groups").InsertOne(ctx, newGroup)

	if err != nil {
		Error.Println(err)
		return false
	}

	return true
}

func (d *MongoDatabase) DeleteGroup(ctx context.Context, id string) bool {
	filter := bson.D{{"id", id}}
	_, err := d.Database.Collection("Groups").DeleteOne(ctx, filter)

	if err != nil {
		Error.Println(err)
		return false
	}

	return true
}
