package Middlewares

import (
	"GSFH/GlobalV"
	"GSFH/Utils"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type GasBody struct {
	GasData string `json:"Gas_Data"`
	GasSign string `json:"Gas_Sign"`
}

const (
	Second = 1
	Minute = 60 * Second
	Hour   = 60 * Minute
)

func ResultResponse(value []byte) string{

	var result GasBody
	cov := string(value)
	conn := GlobalV.RedisGlobalV.Get()

	data := fmt.Sprintf("token_%s_%s",cov, Utils.EncodeMD5(cov))
	res := GetSHA1(data)

	_, err := conn.Do("Set", cov, res)
	if err != nil{
		fmt.Println("Create Response SHA1 Failed")
		return ""
	}
	_, err1 := conn.Do("EXPIRE",cov,Hour)
	if err1 != nil{
		fmt.Println("Create Response EXPIRE TIME Failed")
		return ""
	}

	result.GasData = string(value)
	result.GasSign = res

	cmd1, _ := json.Marshal(result)
	result1 := base64.StdEncoding.EncodeToString(cmd1)

	return result1

}

func GetSHA1(value string) string {
	data := sha1.New()
	data.Write([]byte(value))
	return fmt.Sprintf("%x",data.Sum(nil))
}

func GetGasData() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		 getDataBody := ctx.Request().Body()

			dataBody, dataBodyErr := base64.StdEncoding.DecodeString(string(getDataBody))
			if dataBodyErr != nil {
				fmt.Println("decode error:", dataBodyErr)
				return dataBodyErr
			}

			var StructureGasData *GasBody
			DataUnmarshalErr := json.Unmarshal(dataBody, &StructureGasData)
			if DataUnmarshalErr != nil{
				fmt.Printf("Data Unmarshal Failed %s :  ", DataUnmarshalErr)
				return errors.New("DataUnmarshal Failed")
			}
		CodenameLiquidKey := Utils.EncodeMD5("ThisIsGServerKey")
		DataVerify := hmac.New(sha1.New, []byte(CodenameLiquidKey))
		DataVerify.Write([]byte(StructureGasData.GasData))
		DataVerifyHexDigest := hex.EncodeToString(DataVerify.Sum(nil))

			if StructureGasData.GasSign != DataVerifyHexDigest{
				fmt.Println("Verify Gas GasSign Failed")
				return errors.New("Verify-GasSign Failed")
			}

			DecodeData, DecodeDataErr := base64.StdEncoding.DecodeString(StructureGasData.GasData)
			if DecodeDataErr != nil{
				fmt.Println("DecodeDataErr Err")
				return errors.New("DecodeDataErr Err")
			}

			fmt.Println(string(DecodeData))
			//ctx.Set("CommandData", string(DecodeData))
			ctx.Locals("MiddleData", string(DecodeData))
			ctx.Next()
		return nil
	}

}

