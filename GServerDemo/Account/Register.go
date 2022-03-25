package Account

import (
	"GSFH/Models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type CmdRegister struct {
	FromType string `json:"from_type"`
	Account  string `json:"from_id"`
	Password string `json:"from_token"`
}

func Register(c *fiber.Ctx) error{

	getRegData := c.Locals("MiddleData")
	getRegDataString := getRegData.(string)

	var cmdReg CmdRegister
	getRegDataStringErr := json.Unmarshal([]byte(getRegDataString), &cmdReg)
	if getRegDataStringErr != nil{
		fmt.Println("[Register] getRegDataString Json Unmarshal Failed!")
		return getRegDataStringErr
	}
	toLowPw := strings.ToLower(cmdReg.Password)
	toBase64 := []byte(toLowPw)
	sEnc := base64.StdEncoding.EncodeToString(toBase64)

	res := Models.CreateMember(cmdReg.Account, sEnc)
	if res != nil {
		c.JSON(DataExist)
		return nil
	}

	var mapResult interface{}
	errMap := json.Unmarshal([]byte(getRegDataString), &mapResult)
	if errMap != nil{

	}

	c.JSON(mapResult)
	return nil
}