package config

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewCloudinary(viperConfig *viper.Viper, log *logrus.Logger) *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromParams(
		viperConfig.GetString("CLOUDINARY_CLOUD_NAME"),
		viperConfig.GetString("CLOUDINARY_API_KEY"),
		viperConfig.GetString("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatalf("Failed to creating cloudinary instance : %+v", err)
	}

	return cld
}
