package helper

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Converts uploaded file into buffer
func FormFileToBuffer(log *logrus.Logger, file multipart.File) (*bytes.Buffer, error) {
	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, file); err != nil {
		log.Warnf("Failed to convert file content into buffer : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return buff, nil
}
