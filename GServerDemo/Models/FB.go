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

type FbUser struct {
	AutoId     string    `json:"auto_id" bson:"auto_id"`
	InviteCode string    `json:"invite_code" bson:"invite_code"`
	Create     time.Time `json:"create" bson:"create"`
	Update     time.Time `json:"update" bson:"update"`
	FromType   string    `json:"from_type" bson:"from_type"`
	FromId     string    `json:"from_id" bson:"from_id"`
	FromToken  string    `json:"from_token" bson:"from_token"`
}

type Data struct {
	AppId               string   `json:"app_id" bson:"app_id"`
	Type                string   `json:"type" bson:"type"`
	Application         string   `json:"application" bson:"application"`
	DataAccessExpiresAt int      `json:"data_access_expires_at" bson:"data_access_expires_at"`
	ExpiresAt           int      `json:"expires_at" bson:"expires_at"`
	IsValid             bool     `json:"is_valid" bson:"is_valid"`
	Scopes              []string `json:"scopes" bson:"scopes"`
	UserId              string   `json:"user_id" bson:"user_id"`
}

func CheckFbUser(FbId string) error {

	nowTime := time.Now()
	filter := bson.M{"from_id" : FbId}
	update := bson.M{"$set" : bson.M{"updateTime" : nowTime}}
	var Reg *FbUser
	err := GlobalV.MongoGlobalV.Database("GameDB").Collection("User").FindOneAndUpdate(nil, filter,update).Decode(&Reg)
	if err != nil{
		fmt.Printf("[FB] Can't Find FbID : %s", FbId)
		return err
	}
	return nil
}

func CreateFbUser(FromType, FromId, FbAccessToken string) error {
	
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
		"from_type":   FromType,
		"from_id":     FromId,
		"from_token":  FbAccessToken,
		"createTime":  nowTime,
		"updateTime":  nowTime,
		"invite_code": GetInviteCode,
		"auto_id":     AdminAutoID + InitAutoId + 1,
	}

	user := GlobalV.MongoGlobalV.Database("GameDB").Collection("User")

	_, err := user.Indexes().CreateOne(context.TODO(),mod)
	if err != nil{
		fmt.Println("帳號創建失敗了")
	}

	_, err1 := user.InsertOne(nil, doc)
	if err1 != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

