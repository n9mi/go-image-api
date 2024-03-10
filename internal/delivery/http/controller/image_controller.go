package controller

import (
	"go-image-api/internal/model"
	"go-image-api/internal/usecase"
	"slices"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ImageController struct {
	Log          *logrus.Logger
	ImageUseCase *usecase.ImageUseCase
}

func NewImageController(log *logrus.Logger, imageUseCase *usecase.ImageUseCase) *ImageController {
	return &ImageController{
		Log:          log,
		ImageUseCase: imageUseCase,
	}
}

func (ct *ImageController) ConvertPNGToJPEG(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		ct.Log.Warnf("Failed to parse request multipart/form : %+v", err)
		return fiber.ErrBadRequest
	}

	// Get first uploaded files (if multiple files are uploaded) and only process the first file
	if len(form.File["image"]) == 0 {
		ct.Log.Warn("Validation error : 'image' field is required")
		return fiber.NewError(fiber.StatusBadRequest, "'image' is required")
	}
	file := form.File["image"][0]

	// Validate header, only accepts image/png header
	if file.Header["Content-Type"][0] != "image/png" {
		ct.Log.Warn("Validation error : file header is not image/png")
		return fiber.NewError(fiber.StatusBadRequest, "file should be in png")
	}

	// Send request to usecase
	request := &model.ImageRequest{ImageFileHeader: file}
	response, err := ct.ImageUseCase.ConvertPNGToJPEG(c.UserContext(), request)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *ImageController) Resize(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		ct.Log.Warnf("Failed to parse request multipart/form : %+v", err)
		return fiber.ErrBadRequest
	}

	// Get first uploaded files (if multiple files are uploaded) and only process the first file
	if len(form.File["image"]) == 0 {
		ct.Log.Warn("Validation error : 'image' field is required")
		return fiber.NewError(fiber.StatusBadRequest, "'image' is required")
	}
	file := form.File["image"][0]

	// Validate header, only accepts image/png, image/jpg, image/jpeg header
	extConstraint := []string{
		"image/png",
		"image/jpg",
		"image/jpeg",
	}
	if !slices.Contains(extConstraint, file.Header["Content-Type"][0]) {
		ct.Log.Warn("Validation error : file header is not image/png, image/jpg, or image/jpeg")
		return fiber.NewError(fiber.StatusBadRequest, "file should be in png, jpg, or jpeg")
	}

	// Send request to usecase
	widthReq, _ := strconv.Atoi(c.FormValue("width_in_pixels"))
	heightReq, _ := strconv.Atoi(c.FormValue("height_in_pixels"))
	request := &model.ImageResizeRequest{
		WidthInPixels:   widthReq,
		HeightInPixels:  heightReq,
		ImageFileHeader: file,
	}
	response, err := ct.ImageUseCase.ResizeImage(c.UserContext(), request)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *ImageController) Compress(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		ct.Log.Warnf("Failed to parse request multipart/form : %+v", err)
		return fiber.ErrBadRequest
	}

	// Get first uploaded files (if multiple files are uploaded) and only process the first file
	if len(form.File["image"]) == 0 {
		ct.Log.Warn("Validation error : 'image' field is required")
		return fiber.NewError(fiber.StatusBadRequest, "'image' is required")
	}
	file := form.File["image"][0]

	// Validate header, only accepts image/png, image/jpg, image/jpeg header
	extConstraint := []string{
		"image/png",
		"image/jpg",
		"image/jpeg",
	}
	if !slices.Contains(extConstraint, file.Header["Content-Type"][0]) {
		ct.Log.Warn("Validation error : file header is not image/png, image/jpg, or image/jpeg")
		return fiber.NewError(fiber.StatusBadRequest, "file should be in png, jpg, or jpeg")
	}

	// If 'compress_quality' field is empty, set default as 70
	qualityReq, _ := strconv.Atoi(c.FormValue("compress_quality"))
	if qualityReq < 1 {
		qualityReq = 70
	}

	// Send request to usecase
	request := &model.ImageCompressRequest{
		CompressQuality: qualityReq,
		ImageFileHeader: file,
	}
	response, err := ct.ImageUseCase.CompressImage(c.UserContext(), request)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
