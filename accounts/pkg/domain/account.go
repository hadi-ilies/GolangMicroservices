package domain

import (
	//We will be using mongo that’s why we use a bson.ObjectId
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Auth struct {
	Email    string `json:"email"  bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Account struct {
	// _id allow uniqueness of the field
	Id        bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt time.Time     `bson:"createdat" json:"-"`
	Email     string        `json:"email"  bson:"email"`
	Username  string        `json:"username" bson:"username"`
	Password  string        `json:"password" bson:"password"`
	Balance   uint64        `json:"balance" bson:"balance"`
}
