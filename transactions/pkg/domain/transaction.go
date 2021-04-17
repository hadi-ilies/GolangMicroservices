package domain

import (
	//We will be using mongo that’s why we use a bson.ObjectId
	"gopkg.in/mgo.v2/bson"
)

type Transaction struct {
	// _id allow uniqueness of the field
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Messages []string      `json:"messages"  bson:"messages"`
	Bids     []uint64      `json:"bids" bson:"bids"`
	Status   bool          `json:"status" bson:"status"`
}
