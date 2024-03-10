package config

import (
	"go-image-api/database/migrator"
	"go-image-api/internal/delivery/http/controller"
	"go-image-api/internal/delivery/http/route"
	"go-image-api/internal/repository"
	"go-image-api/internal/usecase"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ConfigBootstrap struct {
	ViperConfig *viper.Viper
	Log         *logrus.Logger
	App         *fiber.App
	DB          *gorm.DB
	Validate    *validator.Validate
	Cloudinary  *cloudinary.Cloudinary
}

func Bootstrap(configBootstrap *ConfigBootstrap) {
	// Setup the repository
	repositorySetup := repository.Setup()

	// Setup the usecase
	useCaseSetup := usecase.Setup(
		configBootstrap.ViperConfig,
		configBootstrap.DB,
		configBootstrap.Validate,
		configBootstrap.Log,
		configBootstrap.Cloudinary,
		repositorySetup,
	)

	// Setup the controller
	controllerSetup := controller.Setup(configBootstrap.Log, useCaseSetup)

	// Setup the routes
	routeConfig := route.RouteConfig{
		App:             configBootstrap.App,
		ControllerSetup: controllerSetup,
	}
	routeConfig.Setup()

	// Drop the database
	if err := migrator.Drop(configBootstrap.DB); err != nil {
		configBootstrap.Log.Fatalf("Failed to drop the database: %+v", err)
	}

	// Migrate the database
	if err := migrator.Migrate(configBootstrap.DB); err != nil {
		configBootstrap.Log.Fatalf("Failed to migrate the database: %+v", err)
	}
}
