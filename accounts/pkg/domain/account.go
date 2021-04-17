package domain

import (
	//We will be using mongo that’s why we use a bson.ObjectId
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	// _id allow uniqueness of the field
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Email    string        `json:"email"  bson:"_id"`
	Username string        `json:"username" bson:"_id"`
	Password string        `json:"password" bson:"password"`
	Balance  uint64        `json:"balance" bson:"balance"`
}
