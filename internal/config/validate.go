package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidate(viperConfig *viper.Viper) *validator.Validate {
	return validator.New()
}
