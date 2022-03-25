package Models

import (
	"GSFH/GlobalV"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"strconv"
	"time"
)

type Member struct {
	Account    string    `json:"account" bson:"account"`
	Password   string    `json:"password" bson:"password"`
}

type User struct {
	AutoId     string    `json:"auto_id" bson:"auto_id"`
	InviteCode string    `json:"invite_code" bson:"invite_code"`
	Create     time.Time `json:"create" bson:"create"`
	Update     time.Time `json:"update" bson:"update"`
	FromType   string    `json:"from_type" bson:"from_type"`
	FromId     string    `json:"from_id" bson:"from_id"`
	FromToken  string    `json:"from_token" bson:"from_token"`
}

func CheckLoginDbMemberData(Account ,Password, collections string) error{
	filter := bson.M{"account" : Account, "password" : Password}
	var user *Member
	err := GlobalV.MongoGlobalV.Database("GameDB").Collection(collections).FindOne(nil, filter).Decode(&user)
	if err != nil{
		fmt.Printf("[Login-CheckLoginDbMemberData] The Login not exist with : %s\n", Account)
		return err
	}
	return nil
}

func CheckLoginDbUserData(FromID ,FromToken, collections string) (error, *User){
	filter := bson.M{"from_id" : FromID, "from_token" : FromToken}
	var user *User
	err := GlobalV.MongoGlobalV.Database("GameDB").Collection(collections).FindOne(nil, filter).Decode(&user)
	if err != nil{
		fmt.Printf("[Login-CheckLoginDbUserData] The Login not exist with : %s\n", FromID)
		return err, nil
	}
	return nil, user
}

func CreateUserLogin(Account, Password, FromType string) (error, *User) {

	mod := mongo.IndexModel{
		Keys: bson.M{
			"from_id": 1,

		}, Options: options.Index().SetUnique(true),
	}

	CreateTime := time.Now()
	member := GlobalV.MongoGlobalV.Database("GameDB").Collection("User")
	doc := bson.M{
		"createTime":  CreateTime,
		"updateTime":  CreateTime,
		"from_id":     Account,
		"from_token":  Password,
		"from_type":   FromType,
	}
	_, err := member.Indexes().CreateOne(context.TODO(),mod)
	if err != nil{
		fmt.Printf("The account number is repeated%s\n", err)
		return err, nil
	}

	_, err2 := member.InsertOne(context.TODO(), doc)
	if err2 != nil {
		fmt.Printf("[Register] Error in insert Dbs2 : %s", err2)
		return err2, nil
	}
	autoIDErr := UpdateLoginUserAutoId(Account)
	if autoIDErr != nil{
		return autoIDErr, nil
	}

	_, userData := CheckLoginDbUserData(Account, Password, "User")
	return nil, userData
}

func UpdateLoginUserAutoId(Account string) error {

	AdminAutoID := UpdateAdminAutoID()
	updateTime := time.Now()

	AutoID2InviteCode := AdminAutoID + InitAutoId + 1
	AutoID64 := int64(AutoID2InviteCode)
	strAutoID64 := strconv.FormatInt(AutoID64, 10)

	GetInviteCode := getAutoIdToInviteCode(strAutoID64)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	autoIdtoString := strconv.Itoa(AdminAutoID + InitAutoId + 1)

	filter := bson.M{"from_id" : Account}
	update := bson.M{"$set" : bson.M{"auto_id" : autoIdtoString, "invite_code" : GetInviteCode,"updateTime" : updateTime}}
	var doc *User
	err := GlobalV.MongoGlobalV.Database("GameDB").Collection("User").FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
		).Decode(&doc)

	if err != nil{
		fmt.Printf("UpdateLoginUserAutoId %s",err)
		return err
	}
	return nil
}

var convertTable = [...]string{
	"e", "d", "c", "b", "a",
	"A", "B", "C", "D", "E",
	"F", "G", "H", "I", "J",
	"j", "i", "h", "g", "f",
	"K", "L", "M", "N", "O",
}

func getAutoIdToInviteCode(autoId string) string {
	inviteCode := ""
	InviteCodeList := make([]int, 0)
	autoIdInt, _ := strconv.ParseInt(autoId, 10, 64)
	rand.Seed(autoIdInt + time.Now().UnixNano())
	newId := strconv.FormatInt(time.Now().Unix()+rand.Int63n(time.Now().Unix()), 10)
	newId = newId + autoId
	newId2Int, _ := strconv.Atoi(newId[len(newId)-13:])
	resultToBase := decimalToBase(InviteCodeList, newId2Int)
	for _, base := range resultToBase {
		inviteCode += convertTable[base]
	}
	return inviteCode
}

func decimalToBase(baseList []int, decimal int) []int {
	base := len(convertTable)
	baseList = append(baseList, decimal%base)
	div := decimal / base
	if div == 0 {
		return baseList
	}
	return decimalToBase(baseList, div)
}
