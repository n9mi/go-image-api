package main

import (
	"fmt"
	"go-image-api/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	app := config.NewFiber(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidate(viperConfig)
	cld := config.NewCloudinary(viperConfig, log)

	configBootstrap := &config.ConfigBootstrap{
		ViperConfig: viperConfig,
		Log:         log,
		App:         app,
		DB:          db,
		Validate:    validate,
		Cloudinary:  cld,
	}
	config.Bootstrap(configBootstrap)

	webPort := viperConfig.GetInt("APP_PORT")
	if err := app.Listen(fmt.Sprintf(":%d", webPort)); err != nil {
		log.Fatalf("Failed to start the app : %+v", err)
	}
}
