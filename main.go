package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	db "pki-server/models"
	"pki-server/routers"
)

func init() {
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("GIN_HOST", ":8080")
	viper.SetDefault("DB_URL", "./pki.sqlite")
	viper.SetDefault("DB_TYPE", "sqlite")
	viper.AutomaticEnv()
	db.ConnectDB(viper.GetString("DB_URL"), viper.GetString("DB_TYPE"))
}

func main() {
	gin.SetMode(viper.GetString("GIN_MODE"))
	r := routers.SetupRouter()
	r.Run(viper.GetString("GIN_HOST"))
}
