package Login

import (
	//"GSFH/Models"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
)

var  (
	appId = "fbAppID"
	appSecret = "fbAppSecret"
	appToken = "fbToken"
)

type AppToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type FBLogicReq struct {
	FromType      string `json:"from_type"`
	FbId          string `json:"from_id"`
	FbAccessToken string `json:"from_token"`
}

type NewData struct {
	Data Data `json:"data"`
}

type Data struct {
	AppId               string   `json:"app_id"`
	Type                string   `json:"type"`
	Application         string   `json:"application"`
	DataAccessExpiresAt int      `json:"data_access_expires_at"`
	ExpiresAt           int      `json:"expires_at"`
	IsValid             bool     `json:"is_valid"`
	Scopes              []string `json:"scopes"`
	UserId              string   `json:"user_id"`
}

func FBIcon(c *fiber.Ctx) error  {
	c.Render("FBLogin", nil)
	return nil
}

func FBLogin(c *fiber.Ctx) error {

	var cmd *FBLogicReq
	err := c.BodyParser(&cmd)
	if err != nil{
		c.JSON(JsonDecodeErrOrDataNull)
		return err
	}
	c.JSON(OK)
	return nil
}


func CheckAppToken() bool {
	appTokUrl := fmt.Sprintf(`https://graph.facebook.com/oauth/access_token?client_id=%s&client_secret=%s&grant_type=client_credentials`, appId, appSecret)
	resp, err := http.Get(appTokUrl)
	if err != nil {
		fmt.Printf("CheckAppToken Fail:%v ", resp)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var atk *AppToken
	data := []byte(string(body))
	json.Unmarshal(data, &atk)
	json.Marshal(atk)

	if atk.AccessToken == ""{
		return false
	}
	return true
}

func CheckAccessToken(userId, accessToken string) bool {
	tokUrl := fmt.Sprintf(`https://graph.facebook.com/debug_token?input_token=%s&access_token=%s`, accessToken, appToken)
	resp, _ := http.Get(tokUrl)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var ack *NewData
	data := []byte(string(body))
	err := json.Unmarshal(data, &ack)
	if err != nil{
		fmt.Println("Get accessToken Fail")
	}

	if ack.Data.IsValid || ack.Data.UserId == userId{
		return true
	}
	return false
}


