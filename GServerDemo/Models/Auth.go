package Models

import (
	"GSFH/GlobalV"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type AuthUser struct {
	AutoId     string    `json:"auto_id" bson:"auto_id"`
	InviteCode string    `json:"invite_code" bson:"invite_code"`
	Create     time.Time `json:"create" bson:"create"`
	Update     time.Time `json:"update" bson:"update"`
	FromType   string    `json:"from_type" bson:"from_type"`
	FromId     string    `json:"from_id" bson:"from_id"`
	FromToken  string    `json:"from_token" bson:"from_token"`
}

func CheckUserByAutoId(AutoId string, InviteCode string) *AuthUser {

	filter := bson.M{"auto_id": AutoId, "invite_code": InviteCode}
	var user *AuthUser

	err := GlobalV.MongoGlobalV.Database("GameDB").Collection("User").FindOne(nil, filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user
}

