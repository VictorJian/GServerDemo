package AU

import (
	"GSFH/GlobalV"
	"GSFH/Middlewares"
	"GSFH/Models"
	"GSFH/Utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

type CmdAuth struct {
	AutoId     string `json:"auto_id"`
	InviteCode string `json:"invite_code"`
	Platform   string `json:"platform"`
}

const (
	Second = 1
	Minute = 60 * Second
	Hour   = 60 * Minute
)

type CmdAuthResponse struct {
	GasId    string `json:"Gas_id" bson:"Gas_id"`
	GasToken string `json:"Gas_token" bson:"Gas_token"`
}

func Auth(c *fiber.Ctx) error{

	getMiddleData := c.Locals("MiddleData")
	dataString := getMiddleData.(string)

	var cmd *CmdAuth
	dataStringErr := json.Unmarshal([]byte(dataString), &cmd)
	if dataStringErr != nil{
		fmt.Println("[Auth]dataStringErr Fail")
		return dataStringErr
	}

	if cmd.Platform == "" {
		platformMain := "PC"
		cmd.Platform = platformMain
	}
	user := Models.CheckUserByAutoId(cmd.AutoId, cmd.InviteCode)
	if user == nil{
		c.JSON(liquidUserErr)
		return errors.New("[auth]CheckUserByAutoId Failed")
	}
	userToken := generateNewToken()
	tokenKey := fmt.Sprintf("token_%s_%s", user.AutoId, cmd.Platform)
	conn := GlobalV.RedisGlobalV.Get()
	defer conn.Close()


	_, getRedisErr := redis.String(conn.Do("Get", tokenKey))
	if getRedisErr != nil {
		_, err := conn.Do("Set", tokenKey, userToken)
		if err != nil{
			fmt.Println("Create user Token Failed")
			return err
		}
		_, err1 := conn.Do("EXPIRE",tokenKey,Hour)
		if err1 != nil{
			fmt.Println("Create EXPIRE TIME Failed")
			return err1
		}
		c.JSON(OK)
		return nil

	}else {

		_, err2 := conn.Do("EXPIRE",tokenKey,Hour)
		if err2 != nil{
			fmt.Println("Create EXPIRE TIME Failed when tokenKey exist")
			return err2
		}

		fmt.Printf("tokenKey TTL update succeeded %s\n", tokenKey)

		cmd, _ := json.Marshal(user)
		sEnc1 := Middlewares.ResultResponse(cmd)
		if sEnc1 == ""{
			fmt.Println("Create User Encryption Failed")
			return errors.New(sEnc1)
		}

		c.JSON(TTLUpdateSucceeded)
	}

	return nil
}

func generateNewToken() string {
	authTime := strconv.Itoa(time.Now().Nanosecond())
	return Utils.EncodeMD5(authTime)
}
