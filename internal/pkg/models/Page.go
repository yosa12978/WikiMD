package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Page struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	LastCommitID string             `json:"last_commit_id" bson:"last_commit_id"`
	Commits      []Commit           `json:"commits" bson:"commits"`
}
