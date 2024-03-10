package route

import (
	"go-image-api/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type RouteConfig struct {
	App             *fiber.App
	ControllerSetup *controller.ControllerSetup
}

func (c *RouteConfig) Setup() {
	route := c.App.Group("/api/v1")

	// Setup basic middleware
	route.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Register controller here
	route.Post("/convert-png-to-jpeg", c.ControllerSetup.ImageController.ConvertPNGToJPEG)
	route.Post("/image-resize", c.ControllerSetup.ImageController.Resize)
	route.Post("/image-compress", c.ControllerSetup.ImageController.Compress)
}
