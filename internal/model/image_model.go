package model

import (
	"mime/multipart"
)

type ImageRequest struct {
	ImageFileHeader *multipart.FileHeader `json:"-" validate:"required"`
}

type ImageResizeRequest struct {
	WidthInPixels   int                   `json:"-" validate:"required,gte=1,lte=3000"`
	HeightInPixels  int                   `json:"-" validate:"required,gte=1,lte=3000"`
	ImageFileHeader *multipart.FileHeader `json:"-" validate:"required"`
}

type ImageCompressRequest struct {
	CompressQuality int                   `json:"-" validate:"required,gte=1,lte=99"`
	ImageFileHeader *multipart.FileHeader `json:"-" validate:"required"`
}

type ImageResponse struct {
	OriginalImageLink string `json:"original_image_link"`
	ResultImageLink   string `json:"result_image_link"`
}
