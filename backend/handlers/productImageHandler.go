package handlers

import (
	db "eshbuket/database"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxImageSize = 10 * 1024 * 1024 // 10MB

var allowedImageExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// getUploadDir returns the filesystem directory where uploads should be stored.
// Priority:
// 1) UPLOAD_DIR env var
// 2) /app/uploads if exists (docker compose volume in this repo)
// 3) ./uploads (local dev)
func getUploadDir() string {
	if v := strings.TrimSpace(os.Getenv("UPLOAD_DIR")); v != "" {
		return v
	}
	if st, err := os.Stat("/app/uploads"); err == nil && st.IsDir() {
		return "/app/uploads"
	}
	return "./uploads"
}

// UploadProductImage handles: POST /api/products/:id/image
// Expects multipart/form-data with field "image".
func UploadProductImage(c *gin.Context) {
	productID := c.Param("id")

	// Ensure product exists
	var exists bool
	if err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)", productID).Scan(&exists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}
	if fileHeader.Size > maxImageSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image too large (max 10MB)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedImageExt[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read image"})
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Failed to close file:", err)
		}
	}()

	head := make([]byte, 512)
	n, _ := file.Read(head)
	ct := http.DetectContentType(head[:n])
	if !strings.HasPrefix(ct, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is not an image"})
		return
	}

	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	} else {
		_ = file.Close()
		file, err = fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read image"})
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Println("Failed to close file:", err)
			}
		}()
	}

	uploadDir := getUploadDir()
	relDir := filepath.Join("products", productID)
	filename := uuid.NewString() + ext
	relPath := filepath.Join(relDir, filename)
	absDir := filepath.Join(uploadDir, relDir)
	absPath := filepath.Join(uploadDir, relPath)

	if err := os.MkdirAll(absDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create directory"})
		return
	}

	out, err := os.Create(absPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println("Failed to close file:", err)
		}
	}()

	lr := io.LimitReader(file, maxImageSize+1)
	written, err := io.Copy(out, lr)
	if err != nil {
		_ = os.Remove(absPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}
	if written > maxImageSize {
		_ = os.Remove(absPath)
		c.JSON(http.StatusBadRequest, gin.H{"error": "image too large (max 10MB)"})
		return
	}

	var prev sql.NullString
	_ = db.DB.QueryRow("SELECT image_path FROM products WHERE id=$1", productID).Scan(&prev)

	_, err = db.DB.Exec("UPDATE products SET image_path=$1 WHERE id=$2", relPath, productID)
	if err != nil {
		_ = os.Remove(absPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	if prev.Valid && prev.String != "" && prev.String != relPath {
		_ = os.Remove(filepath.Join(uploadDir, prev.String))
	}

	c.JSON(http.StatusOK, gin.H{
		"product_id": productID,
		"image_url":  fmt.Sprintf("/api/products/%s/image", productID),
	})
}

// GET /api/products/:id/image
func GetProductImage(c *gin.Context) {
	productID := c.Param("id")

	var rel sql.NullString
	err := db.DB.QueryRow("SELECT image_path FROM products WHERE id=$1", productID).Scan(&rel)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if !rel.Valid || rel.String == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	uploadDir := getUploadDir()
	abs := filepath.Join(uploadDir, rel.String)
	abs = filepath.Clean(abs)

	base := filepath.Clean(uploadDir) + string(os.PathSeparator)
	if !strings.HasPrefix(abs+string(os.PathSeparator), base) && abs != filepath.Clean(uploadDir) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image path"})
		return
	}

	c.File(abs)
}
