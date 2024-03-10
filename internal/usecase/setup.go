package usecase

import (
	"go-image-api/internal/repository"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UseCaseSetup struct {
	ImageUseCase *ImageUseCase
}

func Setup(viperConfig *viper.Viper, db *gorm.DB, validate *validator.Validate, log *logrus.Logger,
	cld *cloudinary.Cloudinary, repositorySetup *repository.RepositorySetup) *UseCaseSetup {
	return &UseCaseSetup{
		ImageUseCase: NewImageUseCase(viperConfig, db, validate, log, cld, repositorySetup.HistoryRepository),
	}
}
