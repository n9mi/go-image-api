# go-image-api-v1
Simple RESTful API with [Fiber](https://github.com/gofiber/fiber) and [gocv](https://github.com/hybridgroup/gocv) to upload and perform image processing such as:
- PNG to JPG conversion 
- Image resizing with specified dimension
- Image compression while maintaining reasonable quality, with modifiable parameter

## Prerequisites
1. [Install gocv locally](https://gocv.io/getting-started)
2. [Create Cloudinary account to get API key](https://cloudinary.com/users/register_free)
3. Running postgreSQL locally or using docker for storing image processing histories.

## /api/v1/convert-png-to-jpeg
Performs png image to jpeg image conversion, but Cloudinary takes jpeg image into jpg, so the 'result_image_link' may in jpg, not jpeg. 
### Header
| Key | Value|
| ------------- | ------------- |
| Content-Type  | multipart/form-data |
### Request
| Key | Value|
| ------------- | ------------- |
| image | [file] |
### Response
| Key | Value|
| ------------- | ------------- |
| original_image_link | https://res.cloudinary.com/... |
| result_image_link | https://res.cloudinary.com/... |

## /api/v1/image-resize
Performs image resizing with provided width and height in pixels. Only accept png, jpg, and jpeg.
### Header
| Key | Value|
| ------------- | ------------- |
| Content-Type  | multipart/form-data |
### Request
| Key | Value|
| ------------- | ------------- |
| image | [file] |
| height_in_pixels | 1-3000 |
| width_in_pixels | 1-3000 |
### Response
| Key | Value|
| ------------- | ------------- |
| original_image_link | https://res.cloudinary.com/... |
| result_image_link | https://res.cloudinary.com/... |

## /api/v1/image-compress
Performs image compression with quality parameter (default as 70). Only accept png, jpg, and jpeg.
### Header
| Key | Value|
| ------------- | ------------- |
| Content-Type  | multipart/form-data |
### Request
| Key | Value|
| ------------- | ------------- |
| image | [file] |
| compress_quality | 0-99 |
### Response
| Key | Value|
| ------------- | ------------- |
| original_image_link | https://res.cloudinary.com/... |
| result_image_link | https://res.cloudinary.com/... |

## TODO
- Accepting images in batches