package domain

import (
	//We will be using mongo that’s why we use a bson.ObjectId
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Transaction struct {
	// _id allow uniqueness of the field
	Id        bson.ObjectId `json:"id" bson:"_id"`
	AccountID bson.ObjectId `json:"account_id" bson:"account_id"`
	AdID      bson.ObjectId `json:"ad_id" bson:"ad_id"`
	Messages  []string      `json:"messages"  bson:"messages"`
	Bids      []uint64      `json:"bids" bson:"bids"`
	CreatedAt time.Time     `bson:"createdat" json:"-"`
	Accepted  bool          `bson:"accepted" json:"accepted"`
	Rejected  bool          `bson:"rejected" json:"rejected"`
}
