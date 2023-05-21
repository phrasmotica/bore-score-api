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

// IsInvitedToGroup implements IDatabase
func (*MongoDatabase) IsInvitedToGroup(ctx context.Context, groupId string, username string) bool {
	panic("unimplemented")
}

// UpdateGroupInvitation implements IDatabase
func (*MongoDatabase) UpdateGroupInvitation(ctx context.Context, newGroupInvitation *models.GroupInvitation) bool {
	panic("unimplemented")
}

// GetGroupInvitation implements IDatabase
func (*MongoDatabase) GetGroupInvitation(ctx context.Context, invitationId string) (bool, *models.GroupInvitation) {
	panic("unimplemented")
}

// AddGroupInvitation implements IDatabase
func (*MongoDatabase) AddGroupInvitation(ctx context.Context, newGroupInvitation *models.GroupInvitation) bool {
	panic("unimplemented")
}

// GetGroupInvitations implements IDatabase
func (*MongoDatabase) GetGroupInvitations(ctx context.Context, username string) (bool, []models.GroupInvitation) {
	panic("unimplemented")
}

// GetGroupInvitationsForGroup implements IDatabase
func (*MongoDatabase) GetGroupInvitationsForGroup(ctx context.Context, groupId string) (bool, []models.GroupInvitation) {
	panic("unimplemented")
}

// GetResultsForGroup implements IDatabase
func (*MongoDatabase) GetResultsForGroup(ctx context.Context, groupId string) (bool, []models.Result) {
	panic("unimplemented")
}

// GetGroupMembershipsForGroup implements IDatabase
func (*MongoDatabase) GetGroupMembershipsForGroup(ctx context.Context, groupId string) (bool, []models.GroupMembership) {
	panic("unimplemented")
}

// GetPlayersInGroup implements IDatabase
func (*MongoDatabase) GetPlayersInGroup(ctx context.Context, groupId string) (bool, []models.Player) {
	panic("unimplemented")
}

// IsInGroup implements IDatabase
func (*MongoDatabase) IsInGroup(ctx context.Context, groupId string, username string) bool {
	panic("unimplemented")
}

// AddGroupMembership implements IDatabase
func (*MongoDatabase) AddGroupMembership(ctx context.Context, newGroupMembership *models.GroupMembership) bool {
	panic("unimplemented")
}

// GetGroupMemberships implements IDatabase
func (*MongoDatabase) GetGroupMemberships(ctx context.Context, username string) (bool, []models.GroupMembership) {
	panic("unimplemented")
}

// UserExistsByEmail implements IDatabase
func (*MongoDatabase) UserExistsByEmail(ctx context.Context, email string) bool {
	panic("unimplemented")
}

// GetResult implements IDatabase
func (*MongoDatabase) GetResult(ctx context.Context, resultId string) (bool, *models.Result) {
	panic("unimplemented")
}

// ResultExists implements IDatabase
func (*MongoDatabase) ResultExists(ctx context.Context, resultId string) bool {
	panic("unimplemented")
}

// GetApprovals implements IDatabase
func (*MongoDatabase) GetApprovals(ctx context.Context, resultId string) (bool, []models.Approval) {
	panic("unimplemented")
}

// AddApproval implements IDatabase
func (*MongoDatabase) AddApproval(ctx context.Context, newApproval *models.Approval) bool {
	panic("unimplemented")
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
