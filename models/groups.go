package models

type Group struct {
	ID             string              `json:"id" bson:"id"`
	Name           string              `json:"name" bson:"name"`
	TimeCreated    int64               `json:"timeCreated" bson:"timeCreated"`
	DisplayName    string              `json:"displayName" bson:"displayName"`
	Description    string              `json:"description" bson:"description"`
	ProfilePicture string              `json:"profilePicture" bson:"profilePicture"`
	CreatedBy      string              `json:"createdBy" bson:"createdBy"`
	Visibility     GroupVisibilityName `json:"visibility" bson:"visibility"`
}

type GroupType struct {
	ID          string              `json:"id" bson:"id"`
	Name        GroupVisibilityName `json:"name" bson:"name"`
	DisplayName string              `json:"displayName" bson:"displayName"`
	Description string              `json:"description" bson:"description"`
}

type GroupVisibilityName string

const (
	Public  GroupVisibilityName = "public"  // players can join whenever they want
	Global  GroupVisibilityName = "global"  // public and everyone is automatically a member
	Private GroupVisibilityName = "private" // players can only join if they are invited
)
