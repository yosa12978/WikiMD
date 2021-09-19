package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Commit struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Body string             `json:"body" bson:"body"`
	Page string             `json:"page" bson:"page"`
	User string             `json:"user" bson:"user"`
	Time int64              `json:"time" bson:"time"`
}
