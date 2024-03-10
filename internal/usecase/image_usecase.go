package usecase

import (
	"bytes"
	"context"
	"go-image-api/internal/entity"
	"go-image-api/internal/helper"
	"go-image-api/internal/model"
	"go-image-api/internal/repository"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"slices"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gocv.io/x/gocv"
	"gorm.io/gorm"
)

type ImageUseCase struct {
	ViperConfig       *viper.Viper
	DB                *gorm.DB
	Validate          *validator.Validate
	Log               *logrus.Logger
	Cloudinary        *cloudinary.Cloudinary
	HistoryRepository *repository.HistoryRepository
}

func NewImageUseCase(viperConfig *viper.Viper, db *gorm.DB, validate *validator.Validate, log *logrus.Logger,
	cld *cloudinary.Cloudinary, historyRepository *repository.HistoryRepository) *ImageUseCase {
	return &ImageUseCase{
		ViperConfig:       viperConfig,
		DB:                db,
		Validate:          validate,
		Log:               log,
		Cloudinary:        cld,
		HistoryRepository: historyRepository,
	}
}

func (u *ImageUseCase) ConvertPNGToJPEG(ctx context.Context, request *model.ImageRequest) (*model.ImageResponse, error) {
	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Validation error : %+v", err)
		return nil, err
	}

	// Open file header
	imageFile, err := request.ImageFileHeader.Open()
	if err != nil {
		u.Log.Warnf("Failed to open file content : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	defer imageFile.Close()

	// Convert file to buffer
	fileBuff, err := helper.FormFileToBuffer(u.Log, imageFile)
	if err != nil {
		u.Log.Warnf("Failed to open file content : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	defer fileBuff.Reset()

	// Validate if file is in png
	originalImageBytes := fileBuff.Bytes()
	if http.DetectContentType(originalImageBytes) != "image/png" {
		u.Log.Warn("Validation error : file is not in png")
		return nil, fiber.NewError(fiber.StatusBadRequest, "file should be in png")
	}

	// Decode image in png
	imagePng, err := png.Decode(bytes.NewReader(originalImageBytes))
	if err != nil {
		u.Log.Warnf("Failed to decode png : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create new buffer in jpeg
	newBuff := new(bytes.Buffer)
	if err := jpeg.Encode(newBuff, imagePng, nil); err != nil {
		u.Log.Warnf("Failed to convert image into jpeg : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	defer newBuff.Reset()

	// Creating history
	ogBuffImage, _, err := image.DecodeConfig(bytes.NewBuffer(fileBuff.Bytes()))
	if err != nil {
		u.Log.Warnf("Failed to decode original image buffer to image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newBuffImage, _, err := image.DecodeConfig(bytes.NewBuffer(newBuff.Bytes()))
	if err != nil {
		u.Log.Warnf("Failed to decode converted image buffer to image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newHistory := &entity.History{
		Timestamp:        time.Now(),
		Type:             "convert_png_jpeg",
		ExtensionBefore:  "image/png",
		ExtensionAfter:   "image/jpg",
		SizeBeforeInMB:   helper.ConvertByteToMB(len(originalImageBytes)),
		SizeAfterInMB:    helper.ConvertByteToMB(len(newBuff.Bytes())),
		HeightBeforeInPx: ogBuffImage.Height,
		WidthBeforeInPx:  ogBuffImage.Width,
		HeightAfterInPx:  newBuffImage.Height,
		WidthAfterInPx:   newBuffImage.Width,
	}

	// Upload original image to cloudinary
	uuidOriginal := uuid.New()
	originalID := "png_" + uuidOriginal.String()
	originalCldResponse, err := u.Cloudinary.Upload.Upload(ctx, fileBuff, uploader.UploadParams{
		PublicID: originalID,
	})
	if err != nil {
		u.Log.Warnf("Failed to upload original image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newHistory.ImageLinkBefore = originalCldResponse.SecureURL

	// Upload converted image to cloudinary
	uuidConverted := uuid.New()
	convertedID := "jpeg_" + uuidConverted.String()
	convertedCldResponse, err := u.Cloudinary.Upload.Upload(ctx, newBuff, uploader.UploadParams{
		PublicID: convertedID,
	})
	if err != nil {
		u.Log.Warnf("Failed to upload converted image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newHistory.ImageLinkAfter = convertedCldResponse.SecureURL

	// Commit history into DB
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := u.HistoryRepository.Repository.Create(tx, newHistory); err != nil {
		u.Log.Warnf("Error adding history : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Error committing history : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	response := &model.ImageResponse{
		OriginalImageLink: originalCldResponse.SecureURL,
		ResultImageLink:   convertedCldResponse.SecureURL,
	}
	return response, nil
}

func (u *ImageUseCase) ResizeImage(ctx context.Context, request *model.ImageResizeRequest) (*model.ImageResponse, error) {
	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Validation error : %+v", err)
		return nil, err
	}

	// Open file header
	imageFile, err := request.ImageFileHeader.Open()
	if err != nil {
		u.Log.Warnf("Failed to open file content : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	defer imageFile.Close()

	// Convert file to buffer
	fileBuff, err := helper.FormFileToBuffer(u.Log, imageFile)
	if err != nil {
		u.Log.Warnf("Failed to open file content : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	defer fileBuff.Reset()

	// Validate if file is in png, jpg, or jpeg
	extConstraint := []string{
		"image/png",
		"image/jpg",
		"image/jpeg",
	}
	originalImageBytes := fileBuff.Bytes()
	contentType := http.DetectContentType(originalImageBytes)
	if !slices.Contains(extConstraint, contentType) {
		u.Log.Warn("Validation error : file is not in png, jpg, or jpeg")
		return nil, fiber.NewError(fiber.StatusBadRequest, "file should be in png, jpg, or jpeg")
	}

	// Convert image bytes to Mat
	originalMat, err := gocv.IMDecode(originalImageBytes, gocv.IMReadAnyColor)
	if err != nil {
		u.Log.Warnf("Failed to convert image to Mat : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Perform resizing
	newMat := gocv.NewMat()
	gocv.Resize(originalMat, &newMat, image.Point{
		X: request.WidthInPixels, Y: request.HeightInPixels}, 0, 0, gocv.InterpolationLinear)

	// Convert Mat into native buffer
	var newBuffNative *gocv.NativeByteBuffer
	switch contentType {
	case "image/png":
		newBuffNative, err = gocv.IMEncode(".png", newMat)
	case "image/jpg":
		newBuffNative, err = gocv.IMEncode(".jpg", newMat)
	case "image/jpeg":
		newBuffNative, err = gocv.IMEncode(".jpeg", newMat)
	}
	if err != nil {
		u.Log.Warnf("Failed to convert Mat to buffer : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Convert gocv native buffer to buffer
	newBuff := bytes.NewBuffer(newBuffNative.GetBytes())
	defer newBuff.Reset()

	// Creating history
	ogBuffImage, _, err := image.DecodeConfig(bytes.NewBuffer(fileBuff.Bytes()))
	if err != nil {
		u.Log.Warnf("Failed to decode original image buffer to image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newBuffImage, _, err := image.DecodeConfig(bytes.NewBuffer(newBuff.Bytes()))
	if err != nil {
		u.Log.Warnf("Failed to decode converted image buffer to image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newHistory := &entity.History{
		Timestamp:        time.Now(),
		Type:             "resize_image",
		ExtensionBefore:  contentType,
		ExtensionAfter:   contentType,
		SizeBeforeInMB:   helper.ConvertByteToMB(len(originalImageBytes)),
		SizeAfterInMB:    helper.ConvertByteToMB(len(newBuff.Bytes())),
		HeightBeforeInPx: ogBuffImage.Height,
		WidthBeforeInPx:  ogBuffImage.Width,
		HeightAfterInPx:  newBuffImage.Height,
		WidthAfterInPx:   newBuffImage.Width,
	}

	// Upload original image to cloudinary
	uuidOriginal := uuid.New()
	originalID := "original_" + uuidOriginal.String()
	originalCldResponse, err := u.Cloudinary.Upload.Upload(ctx, fileBuff, uploader.UploadParams{
		PublicID: originalID,
	})
	if err != nil {
		u.Log.Warnf("Failed to upload original image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newHistory.ImageLinkBefore = originalCldResponse.SecureURL

	// Upload converted image to cloudinary
	uuidConverted := uuid.New()
	convertedID := "resized_" + uuidConverted.String()
	convertedCldResponse, err := u.Cloudinary.Upload.Upload(ctx, newBuff, uploader.UploadParams{
		PublicID: convertedID,
	})
	if err != nil {
		u.Log.Warnf("Failed to upload converted image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newHistory.ImageLinkAfter = convertedCldResponse.SecureURL

	// Commit history into DB
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := u.HistoryRepository.Repository.Create(tx, newHistory); err != nil {
		u.Log.Warnf("Error adding history : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Error committing history : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	response := &model.ImageResponse{
		OriginalImageLink: originalCldResponse.SecureURL,
		ResultImageLink:   convertedCldResponse.SecureURL,
	}
	return response, nil
}

func (u *ImageUseCase) CompressImage(ctx context.Context, request *model.ImageCompressRequest) (*model.ImageResponse, error) {
	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Validation error : %+v", err)
		return nil, err
	}

	// Open file header
	imageFile, err := request.ImageFileHeader.Open()
	if err != nil {
		u.Log.Warnf("Failed to open file content : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	defer imageFile.Close()

	// Convert file to buffer
	fileBuff, err := helper.FormFileToBuffer(u.Log, imageFile)
	if err != nil {
		u.Log.Warnf("Failed to open file content : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	defer fileBuff.Reset()

	// Validate if file is in png, jpg, or jpeg
	extConstraint := []string{
		"image/png",
		"image/jpg",
		"image/jpeg",
	}
	originalImageBytes := fileBuff.Bytes()
	contentType := http.DetectContentType(originalImageBytes)
	if !slices.Contains(extConstraint, contentType) {
		u.Log.Warn("Validation error : file is not in png, jpg, or jpeg")
		return nil, fiber.NewError(fiber.StatusBadRequest, "file should be in png, jpg, or jpeg")
	}

	// Convert image bytes to Mat
	originalMat, err := gocv.IMDecode(originalImageBytes, gocv.IMReadAnyColor)
	if err != nil {
		u.Log.Warnf("Failed to convert image to Mat : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Petform compression
	var newBuffNative *gocv.NativeByteBuffer
	switch contentType {
	case "image/png":
		newBuffNative, err = gocv.IMEncodeWithParams(".png", originalMat, []int{
			gocv.IMWriteJpegQuality, request.CompressQuality})
	case "image/jpg":
		newBuffNative, err = gocv.IMEncodeWithParams(".jpg", originalMat, []int{
			gocv.IMWriteJpegQuality, request.CompressQuality})
	case "image/jpeg":
		newBuffNative, err = gocv.IMEncodeWithParams(".jpeg", originalMat, []int{
			gocv.IMWriteJpegQuality, request.CompressQuality})
	}
	if err != nil {
		u.Log.Warnf("Failed to convert Mat to buffer : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Convert gocv native buffer to buffer
	newBuff := bytes.NewBuffer(newBuffNative.GetBytes())
	defer newBuff.Reset()

	// Creating history
	ogBuffImage, _, err := image.DecodeConfig(bytes.NewBuffer(fileBuff.Bytes()))
	if err != nil {
		u.Log.Warnf("Failed to decode original image buffer to image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newBuffImage, _, err := image.DecodeConfig(bytes.NewBuffer(newBuff.Bytes()))
	if err != nil {
		u.Log.Warnf("Failed to decode converted image buffer to image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newHistory := &entity.History{
		Timestamp:        time.Now(),
		Type:             "compress_image",
		ExtensionBefore:  contentType,
		ExtensionAfter:   contentType,
		SizeBeforeInMB:   helper.ConvertByteToMB(len(originalImageBytes)),
		SizeAfterInMB:    helper.ConvertByteToMB(len(newBuff.Bytes())),
		HeightBeforeInPx: ogBuffImage.Height,
		WidthBeforeInPx:  ogBuffImage.Width,
		HeightAfterInPx:  newBuffImage.Height,
		WidthAfterInPx:   newBuffImage.Width,
	}

	// Upload original image to cloudinary
	uuidOriginal := uuid.New()
	originalID := "original_" + uuidOriginal.String()
	originalCldResponse, err := u.Cloudinary.Upload.Upload(ctx, fileBuff, uploader.UploadParams{
		PublicID: originalID,
	})
	if err != nil {
		u.Log.Warnf("Failed to upload original image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newHistory.ImageLinkBefore = originalCldResponse.SecureURL

	// Upload converted image to cloudinary
	uuidConverted := uuid.New()
	convertedID := "compressed_" + uuidConverted.String()
	convertedCldResponse, err := u.Cloudinary.Upload.Upload(ctx, newBuff, uploader.UploadParams{
		PublicID: convertedID,
	})
	if err != nil {
		u.Log.Warnf("Failed to upload converted image : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	newHistory.ImageLinkAfter = convertedCldResponse.SecureURL

	// Commit history into DB
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := u.HistoryRepository.Repository.Create(tx, newHistory); err != nil {
		u.Log.Warnf("Error adding history : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Error committing history : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	response := &model.ImageResponse{
		OriginalImageLink: originalCldResponse.SecureURL,
		ResultImageLink:   convertedCldResponse.SecureURL,
	}
	return response, nil
}
