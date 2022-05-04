package main

import (
	"sk-integrated-services/pkg/config"
	"sk-integrated-services/pkg/database"
	"sk-integrated-services/pkg/logger"
	"sk-integrated-services/pkg/routers"

	"github.com/spf13/viper"
)

func main() {
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config.SetupConfig() error: %s", err)
	}

	if err := database.Connection(); err != nil {
		logger.Fatalf("database.DbConnection error: %s", err)
	}

	db := database.GetDB()
	router := routers.Routes(db)

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8000")

	logger.Fatalf("%v:%v", router.Run(viper.GetString("server.host")+":"+viper.GetString("server.port")))

}
