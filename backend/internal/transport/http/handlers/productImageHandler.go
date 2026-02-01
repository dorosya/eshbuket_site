package handlers

import (
	"net/http"

	"eshbuket/internal/service/productimage"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	svc *productimage.Service
}

func NewImageHandler(svc *productimage.Service) *ImageHandler {
	return &ImageHandler{svc: svc}
}

// POST /api/products/:id/image
// Expects multipart/form-data with field "image".
func (h *ImageHandler) UploadProductImage(c *gin.Context) {
	productID := c.Param("id")

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read image"})
		return
	}
	defer func() { _ = f.Close() }()

	url, err := h.svc.Upload(c.Request.Context(), productID, fileHeader.Filename, fileHeader.Size, f)
	if err != nil {
		switch err {
		case productimage.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		case productimage.ErrImageTooLarge:
			c.JSON(http.StatusBadRequest, gin.H{"error": "image too large (max 10MB)"})
		case productimage.ErrBadImage:
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported or invalid image"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product_id": productID,
		"image_url":  url,
	})
}

// GET /api/products/:id/image
func (h *ImageHandler) GetProductImage(c *gin.Context) {
	productID := c.Param("id")

	f, name, mod, ct, etag, err := h.svc.Open(c.Request.Context(), productID)
	if err != nil {
		switch err {
		case productimage.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		case productimage.ErrImageNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		case productimage.ErrBadImage:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image path"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open image"})
		}
		return
	}
	defer func() { _ = f.Close() }()

	if ct != "" {
		c.Header("Content-Type", ct)
	}
	c.Header("Cache-Control", "public, max-age=3600")
	if etag != "" {
		c.Header("ETag", etag)
	}

	http.ServeContent(c.Writer, c.Request, name, mod, f)
}
