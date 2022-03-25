package Login

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
)

type GoogleLogicReq struct {
	FromType  string `json:"from_type"`
	FromId    string `json:"from_id"`
	FromToken string `json:"from_token"`
}

type GoogleData struct {
	ISS           string `json:"iss"`
	AZP           string `json:"azp"`
	AUD           string `json:"aud"`
	SUB           string `json:"sub"`
	HD            string `json:"hd"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	IAT           string `json:"iat"`
	EXP           string `json:"exp"`
	JTI           string `json:"jti"`
	ALG           string `json:"alg"`
	KID           string `json:"kid"`
	Type          string `json:"typ"`
}

func GoogleIcon(c *fiber.Ctx) error {

	c.Render("GoogleLogin", nil)
	return nil
}

func GoogleLogin(c *fiber.Ctx) error {

	var cmd *GoogleLogicReq
	err := c.BodyParser(&cmd)
	if err != nil {
		c.JSON(JsonDecodeErrOrDataNull)
		return err
	}
	c.JSON(OK)
	return nil

}

func CheckGoogleIDToken(userID string) bool {
	idTokUrl := fmt.Sprintf(`https://oauth2.googleapis.com/tokeninfo?id_token=%s`, userID)
	resp, err := http.Get(idTokUrl)
	if err != nil {
		fmt.Printf("CheckIDToken Fail:%v ", resp)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var googleData *GoogleData
	data := []byte(string(body))
	unmarshalGoogleErr := json.Unmarshal(data, &googleData)
	if unmarshalGoogleErr != nil {
		fmt.Printf("unmarshalGoogleErr Fail:%v ", unmarshalGoogleErr)
		return false
	}

	_, marshalGoogleErr := json.Marshal(googleData)
	if marshalGoogleErr != nil {
		fmt.Printf("unmarshalGoogleErr Fail:%v ", marshalGoogleErr)
		return false
	}

	if googleData.SUB == "" {
		return false
	}
	return true
}


