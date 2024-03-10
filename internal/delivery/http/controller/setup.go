package controller

import (
	"go-image-api/internal/usecase"

	"github.com/sirupsen/logrus"
)

type ControllerSetup struct {
	ImageController *ImageController
}

func Setup(log *logrus.Logger, useCaseSetup *usecase.UseCaseSetup) *ControllerSetup {
	return &ControllerSetup{
		ImageController: NewImageController(log, useCaseSetup.ImageUseCase),
	}
}
