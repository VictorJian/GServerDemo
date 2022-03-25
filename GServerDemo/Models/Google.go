package Models

import (
	"GSFH/GlobalV"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type GoogleUser struct {
	AutoId     string    `json:"auto_id" bson:"auto_id"`
	InviteCode string    `json:"invite_code" bson:"invite_code"`
	Create     time.Time `json:"create" bson:"create"`
	Update     time.Time `json:"update" bson:"update"`
	FromType   string    `json:"from_type" bson:"from_type"`
	FromId     string    `json:"from_id" bson:"from_id"`
	FromToken  string    `json:"from_token" bson:"from_token"`
}

type GoogleData struct {
	ISS           string `json:"iss" bson:"iss"`
	AZP           string `json:"azp" bson:"azp"`
	AUD           string `json:"aud" bson:"aud"`
	SUB           string `json:"sub" bson:"sub"`
	HD            string `json:"hd" bson:"hd"`
	Email         string `json:"email" bson:"email"`
	EmailVerified bool   `json:"email_verified" bson:"email_verified"`
	AtHash        string `json:"at_hash" bson:"at_hash"`
	Name          string `json:"name" bson:"name"`
	Picture       string `json:"picture" bson:"picture"`
	GivenName     string `json:"given_name" bson:"given_name"`
	FamilyName    string `json:"family_name" bson:"family_name"`
	Locale        string `json:"locale" bson:"locale"`
	IAT           string `json:"iat" bson:"iat"`
	EXP           string `json:"exp" bson:"exp"`
	JTI           string `json:"jti" bson:"jti"`
	ALG           string `json:"alg" bson:"alg"`
	KID           string `json:"kid" bson:"kid"`
	Type          string `json:"typ" bson:"typ"`
}

func CheckGoogleUser(GoogleId string) error {
	nowTime := time.Now()
	filter := bson.M{"from_id" : GoogleId}
	update := bson.M{"$set" : bson.M{"updateTime" : nowTime}}
	var Reg *GoogleUser
	err := GlobalV.MongoGlobalV.Database("GameDB").Collection("User").FindOneAndUpdate(nil, filter,update).Decode(&Reg)
	if err != nil{
		fmt.Printf("[FB] Can't Find FbID : %s", GoogleId)
		return err
	}
	return nil
}

func CreateGoogleUser(fromType, fromId, googleIDToken string) error {

	mod := mongo.IndexModel{
		Keys: bson.M{
			"from_id": 1,

		}, Options: options.Index().SetUnique(true),
	}

	nowTime := time.Now()
	AdminAutoID := UpdateAdminAutoID()

	AutoID2InviteCode := AdminAutoID + InitAutoId + 1
	AutoID64 := int64(AutoID2InviteCode)
	strAutoID64 := strconv.FormatInt(AutoID64, 10)

	GetInviteCode := getAutoIdToInviteCode(strAutoID64)

	doc := bson.M{
		"from_type":   fromType,
		"from_id":     fromId,
		"from_token":  googleIDToken,
		"createTime":  nowTime,
		"updateTime":  nowTime,
		"invite_code": GetInviteCode,
		"auto_id":     AdminAutoID + InitAutoId + 1,
	}

	user := GlobalV.MongoGlobalV.Database("GameDB").Collection("User")

	_, err := user.Indexes().CreateOne(context.TODO(),mod)
	if err != nil{
		fmt.Println("The account creation failed")
	}

	_, err1 := user.InsertOne(nil, doc)
	if err1 != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
