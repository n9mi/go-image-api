package config

import (
	"fmt"
	"go-image-api/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(viperConfig *viper.Viper) *fiber.App {
	bodyLimit := viperConfig.GetInt("BODY_LIMIT_IN_MB")
	if bodyLimit < 1 { // If body limit is not configured, set defaults to 5 MB
		bodyLimit = 5
	}

	app := fiber.New(fiber.Config{
		AppName:      viperConfig.GetString("APP_NAME"),
		ErrorHandler: customErrorHandler(),
		BodyLimit:    bodyLimit * 1024 * 1024,
	})

	return app
}

func customErrorHandler() func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		var response model.ErrorResponse

		if err != nil {
			if errConv, ok := err.(validator.ValidationErrors); ok {
				response.Code = fiber.StatusBadRequest

				for _, errItem := range errConv {
					switch errItem.Tag() {
					case "required":
						response.Messages = append(response.Messages, fmt.Sprintf("%s is required", errItem.Field()))
					case "min":
						response.Messages = append(response.Messages, fmt.Sprintf("%s is should more than %d",
							errItem.Tag(), errItem.Value()))
					case "max":
						response.Messages = append(response.Messages, fmt.Sprintf("%s is should me less than %d",
							errItem.Tag(), errItem.Value()))
					}
				}
			} else if errConv, ok := err.(*fiber.Error); ok {
				response.Code = errConv.Code
				response.Messages = []string{errConv.Message}
			} else {
				response.Code = fiber.StatusInternalServerError
				response.Messages = []string{"Internal server error"}
			}

			return c.Status(response.Code).JSON(response)
		}

		return nil
	}
}
