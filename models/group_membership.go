package models

type GroupMembership struct {
	ID           string `json:"id" bson:"id"`
	GroupID      string `json:"groupId" bson:"groupId"`
	TimeCreated  int64  `json:"timeCreated" bson:"timeCreated"`
	Username     string `json:"username" bson:"username"`
	InvitationID string `json:"invitationId" bson:"invitationId"`
}

type GroupInvitation struct {
	ID               string           `json:"id" bson:"id"`
	GroupID          string           `json:"groupId" bson:"groupId"`
	TimeCreated      int64            `json:"timeCreated" bson:"timeCreated"`
	Username         string           `json:"username" bson:"username"`
	InviterUsername  string           `json:"inviterUsername" bson:"inviterUsername"`
	InvitationStatus InvitationStatus `json:"invitationStatus" bson:"invitationStatus"`
}

type InvitationStatus string

const (
	Sent     InvitationStatus = "sent"
	Accepted InvitationStatus = "accepted"
	Declined InvitationStatus = "declined"
)
