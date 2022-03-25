package Models

import (
	"GSFH/GlobalV"
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

)

type DocUser struct {
	AutoId     int    `json:"auto_id" bson:"auto_id"`
	InviteCode string    `json:"invite_code" bson:"invite_code"`
	Create     time.Time `json:"create" bson:"create"`
	Update     time.Time `json:"update" bson:"update"`
	FromType   string    `json:"from_type" bson:"from_type"`
	FromId     string    `json:"from_id" bson:"from_id"`
	FromToken  string    `json:"from_token" bson:"from_token"`
}

type LiquidAdminConfCounterSetting struct {
	Admin   string `json:"admin" bson:"admin"`
	Counter int    `json:"counter" bson:"counter"`
}

const nitAutoIdTest = 1000000

func ConnectMongo(url string) *mongo.Client {

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client
}



func RedisInit(redisPath string) *redis.Pool  {

	client := &redis.Pool{
		MaxIdle: 8,
		MaxActive: 0,
		IdleTimeout: 3600,
		Dial: func() (redis.Conn, error) {
			res, err := redis.Dial("tcp", redisPath)
			if err != nil{
				return nil, err
			}
			return res, err
		},
	}
	fmt.Println("Connected to Redis!")
	return client
}


func CreateAdminAutoId(){

	admin := GlobalV.MongoGlobalV.Database("GameDB").Collection("Admin")
	doc := bson.M{
		"admin": "auto_id",
		"counter": 0,
	}
	_, err := admin.InsertOne(context.TODO(), doc)
	if err != nil{
		fmt.Println("Insert counter fail")
	}
	return
}
func FindAdminAutoID() int{

	filter := bson.M{}
	findOption := options.Find()
	projection := bson.M{
		"counter": true,
	}
	findOption.SetProjection(projection)

	doc, err := GlobalV.MongoGlobalV.Database("GameDB").Collection("Admin").Find(
		nil,
		filter,
		findOption,
	)

	if err != nil {
		return -1
	}

	var counterList []AdminConfCounterSetting

	err = doc.All(context.TODO(), &counterList)
	if len(counterList) == 0{
		return 1
	}
	return -1
}


