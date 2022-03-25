package Login

import (
	"GSFH/Models"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type CmdLogin struct {
	FromType  string `json:"from_type"`
	FromId    string `json:"from_id"`
	FromToken string `json:"from_token"`
}

type CmdLoginResp struct {
	AutoId     string `json:"auto_id" bson:"auto_id"`
	InviteCode string `json:"invite_code" bson:"invite_code"`
}

func Login(c *fiber.Ctx) error {

	getData := c.Locals("MiddleData")
	dataString := getData.(string)
	var cmd *CmdLogin
	errJ := json.Unmarshal([]byte(dataString), &cmd)
	if errJ != nil{
		fmt.Println(errJ)
		return errJ
	}

	switch cmd.FromType {
	case "VSystem":

		toLowPw := strings.ToLower(cmd.FromToken)
		toBase64 := []byte(toLowPw)
		sEnc := base64.StdEncoding.EncodeToString(toBase64)
		//BASE 可逆，找尋不可逆的

		memberData := Models.CheckLoginDbMemberData(cmd.FromId, sEnc, "Member")
		if memberData != nil{
			c.JSON(DataMismatch)
			return memberData
		}
		userData, _ := Models.CheckLoginDbUserData(cmd.FromId,sEnc,"User")
		if userData != nil{
			data1, loginUser := Models.CreateUserLogin(cmd.FromId,sEnc,"VSystem")
			if data1 != nil{
				c.JSON(CreateUserErr)
				return data1
			}
				mapResult := &CmdLoginResp{
				AutoId: loginUser.AutoId,
				InviteCode: loginUser.InviteCode,
			}

				c.JSON(mapResult)
				return nil
		}

		var mapResult interface{}
		errMap := json.Unmarshal([]byte(dataString), &mapResult)
		if errMap != nil{

		}
		c.JSON(mapResult)
		return nil

	case "FaceBook":

		//Verify the validity of the APP Token
		isAppToken := CheckAppToken()
		if !isAppToken {
			c.JSON(AppTokenErr)
			return nil
		}
		c.JSON(CheckAppTokenOK)

		// Confirm the validity of access token with FB
		isAccess := CheckAccessToken(cmd.FromId, cmd.FromToken)
		if !isAccess {
			c.JSON(AccessTokenErr)
			return errors.New("password Error")
		}
		c.JSON(CheckFBAccessTokenOK)

		checkErr := Models.CheckFbUser(cmd.FromId)
		if checkErr != nil {
			result := Models.CreateFbUser(cmd.FromType, cmd.FromId, cmd.FromToken)
			if result != nil{
				c.JSON(CreateFbUserErr)
				return result
			}
		}

		c.JSON(OK)
		return nil

	case "Google":
		isAccess := CheckGoogleIDToken(cmd.FromToken)
		if !isAccess {
			c.JSON(AccessTokenErr)
			return errors.New("badddddddd")
		}
		c.JSON(CheckGoogleIDTokenOK)

		checkErr := Models.CheckGoogleUser(cmd.FromId)
		if checkErr != nil {
			result := Models.CreateGoogleUser(cmd.FromType, cmd.FromId, cmd.FromToken)
			if result != nil {
				c.JSON(CreateGoogleUserErr)
				return result
			}
		}
		c.JSON(OK)

		var mapResult map[string]interface{}
		errMap := json.Unmarshal([]byte(dataString), &mapResult)
		if errMap != nil{

		}
		c.JSON(mapResult)
		return nil

	default:
		c.JSON(FromTypeErr)

	}
	return nil
}




