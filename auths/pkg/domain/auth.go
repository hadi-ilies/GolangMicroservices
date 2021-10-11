package domain

//We will be using mongo that’s why we use a bson.ObjectId
type Auth struct {
	Email        string `json:"email" bson:"email"`
	Password     string `json:"password" bson:"password"`
	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`
	AccessUuid   string `bson:"access_uuid"`
	RefreshUuid  string `bson:"refresh_uuid"`
	AtExpires    int64  `bson:"at_expires"`
	RtExpires    int64  `bson:"rt_expires"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	AccessUuid   string `json:"access_uuid" bson:"access_uuid"`
	RefreshUuid  string `json:"refresh_uuid" bson:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires" bson:"at_expires"`
	RtExpires    int64  `json:"rt_expires" bson:"rt_expires"`
}

type AccessDetails struct {
	AccessUuid string `json:"access_uuid" bson:"access_uuid"`
	UserID     string `json:"user_id" bson:"user_id"`
}
