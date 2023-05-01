package models

type Approval struct {
	ID             string         `json:"id" bson:"id"`
	ResultID       string         `json:"resultId" bson:"resultId"`
	TimeCreated    int64          `json:"timeCreated" bson:"timeCreated"`
	Username       string         `json:"username" bson:"username"`
	ApprovalStatus ApprovalStatus `json:"approvalStatus" bson:"approvalStatus"`
}

type ApprovalStatus string

const (
	Pending  ApprovalStatus = "pending"
	Approved ApprovalStatus = "approved"
	Rejected ApprovalStatus = "rejected"
)
