package Models

import (
	"GSFH/GlobalV"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CmdRegister struct {
	FromType string `json:"from_type" bson:"from_type"`
	Account  string `json:"account" bson:"account"`
	Password string `json:"password" bson:"password"`
}

type AdminConfCounterSetting struct {
	Admin   string `json:"admin" bson:"admin"`
	Counter int    `json:"counter" bson:"counter"`
}

const InitAutoId = 1000000

func CreateMember(Account, Password string) error {

		mod := mongo.IndexModel{
			Keys: bson.M{
				"account": 1,

			}, Options: options.Index().SetUnique(true),
		}

		CreateTime := time.Now()
		member := GlobalV.MongoGlobalV.Database("GameDB").Collection("Member")
		doc := bson.M{
			"account":    Account,
			"password":   Password,
			"createTime": CreateTime,
		}

		ind, err := member.Indexes().CreateOne(context.TODO(),mod)
		if err != nil{
			fmt.Println("CreateOne() index failed:", ind)
		}

		_, Err := member.InsertOne(context.TODO(), doc)
		if Err != nil {
			fmt.Printf("[Register] Error in insert Dbs : %s", Err)
			return Err
		}
	return nil
}

func UpdateAdminAutoID() int{

	var findAdminAutoID *AdminConfCounterSetting
	filter := bson.M{"admin": "auto_id"}
	update := bson.M{"$inc" : bson.M{"counter" : 1}}
	doc := GlobalV.MongoGlobalV.Database("GameDB").Collection("Admin").FindOneAndUpdate(
		nil,
		filter,
		update,
	).Decode(&findAdminAutoID)

	if doc != nil{
		return -1
	}

	return findAdminAutoID.Counter
}