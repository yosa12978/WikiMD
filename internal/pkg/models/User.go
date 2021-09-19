package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"-" bson:"password"`
	Email    string             `json:"email" bson:"email"`
	Token    string             `json:"token" bson:"token"`
	Role     Role               `json:"role" bson:"role"`
	Regdate  int64              `json:"regdate" bson:"regdate"`
}
