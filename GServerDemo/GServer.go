package main

import (
	"GSFH/AU"
	"GSFH/Account"
	"GSFH/GlobalV"
	"GSFH/Login"
	"GSFH/Middlewares"
	"GSFH/Models"
	"GSFH/Setting"
	"GSFH/Utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	MongoPath string
	RedisPath string
}

const GetKey = "ThisIsGServerKey"

func main() {
//Gserver environment parameter setting
	var config Config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	//viper.SetDefault("Config","")
	err := viper.ReadInConfig()
	if err != nil{
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viperUnmarshalErr := viper.Unmarshal(&config)
	if viperUnmarshalErr != nil{
		panic(fmt.Errorf("Fatal error viperUnmarshalErr: %s \n", viperUnmarshalErr))
	}

//Gserver Dbs environment Settings

	GlobalV.RedisGlobalV = Models.RedisInit(config.RedisPath)
	//GlobalV.MongoGlobalV = Models.ConnectMongo("mongodb://localhost:27017")
	GlobalV.MongoGlobalV = Models.ConnectMongo(config.MongoPath)
	Setting.SettingAdmin()

	vEngine := html.New(".", ".html")
	vEngine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: vEngine,
	})

	app.Use(cors.New())
	app.Get("/@", func(c *fiber.Ctx) error {
		c.SendString(Utils.EncodeMD5(GetKey))
		return nil
	})
	app.Use(Middlewares.GetGasData())
	{
		app.Post("/register", Account.Register)
		app.Post("/login", Login.Login)
		app.Post("/fbLogin", Login.FBLogin)
		app.Post("/googleLogin", Login.GoogleLogin)
		app.Post("/auth", AU.Auth)
	}

	//ThirdPartLogin
	app.Get("/fbIcon", Login.FBIcon)
	app.Get("/googleIcon", Login.GoogleIcon)
	app.Get("/twitterIcon", Login.TwitterIcon)

	//app.Listen(":8899")
	app.Listen(config.ServerPort)

}

