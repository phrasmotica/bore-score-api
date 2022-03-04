package models

type Group struct {
	ID          string        `json:"id" bson:"id"`
	Name        string        `json:"name" bson:"name"`
	TimeCreated int64         `json:"timeCreated" bson:"timeCreated"`
	Type        GroupTypeName `json:"type" bson:"type"`
	DisplayName string        `json:"displayName" bson:"displayName"`
	Description string        `json:"description" bson:"description"`
}

type GroupType struct {
	ID          string        `json:"id" bson:"id"`
	Name        GroupTypeName `json:"name" bson:"name"`
	DisplayName string        `json:"displayName" bson:"displayName"`
	Description string        `json:"description" bson:"description"`
}

type GroupTypeName string

const (
	Public  GroupTypeName = "public"  // players can join whenever they want
	Global  GroupTypeName = "global"  // public and everyone is automatically a member
	Private GroupTypeName = "private" // players can only join if they are invited
)
