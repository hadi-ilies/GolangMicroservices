package domain

import (
	//We will be using mongo that’s why we use a bson.ObjectId
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Ad struct {
	// _id allow uniqueness of the field
	Id         bson.ObjectId `json:"id" bson:"_id"`
	AccountID  bson.ObjectId `json:"-" bson:"account_id"`
	CreatedAt  time.Time     `json:"-" bson:"createdat"`
	Title      string        `json:"title"  bson:"title"`
	Descripion string        `json:"description" bson:"description"`
	Price      uint64        `json:"price" bson:"price"`
	Picture    string        `json:"picture" bson:"picture"`
}
