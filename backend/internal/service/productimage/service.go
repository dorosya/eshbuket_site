package productimage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const MaxImageSize = 10 * 1024 * 1024 // 10MB

var allowedImageExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrImageNotFound   = errors.New("image not found")
	ErrBadImage        = errors.New("invalid image")
	ErrImageTooLarge   = errors.New("image too large")
)

type Repository interface {
	ProductExists(ctx context.Context, productID string) (bool, error)
	GetImagePath(ctx context.Context, productID string) (string, bool, error)
	SetImagePath(ctx context.Context, productID string, relPath string) error
}

type Service struct {
	repo      Repository
	uploadDir string
}

func NewService(repo Repository, uploadDir string) *Service {
	if strings.TrimSpace(uploadDir) == "" {
		uploadDir = defaultUploadDir()
	}
	return &Service{repo: repo, uploadDir: uploadDir}
}

// defaultUploadDir возвращает директорию в которую будут загружаться изображения
// Приоритет:
// 1) UPLOAD_DIR env var
// 2) /app/uploads if exists (docker compose volume in this repo)
// 3) ./uploads (local dev)
func defaultUploadDir() string {
	if v := strings.TrimSpace(os.Getenv("UPLOAD_DIR")); v != "" {
		return v
	}
	if st, err := os.Stat("/app/uploads"); err == nil && st.IsDir() {
		return "/app/uploads"
	}
	return "./uploads"
}

func (s *Service) Upload(ctx context.Context, productID string, originalFilename string, size int64, file io.ReadSeeker) (imageURL string, err error) {
	exists, err := s.repo.ProductExists(ctx, productID)
	if err != nil {
		return "", fmt.Errorf("check product exists: %w", err)
	}
	if !exists {
		return "", ErrProductNotFound
	}

	if size <= 0 {
		return "", ErrBadImage
	}
	if size > MaxImageSize {
		return "", ErrImageTooLarge
	}

	ext := strings.ToLower(filepath.Ext(originalFilename))
	if !allowedImageExt[ext] {
		return "", ErrBadImage
	}

	head := make([]byte, 512)
	n, _ := file.Read(head)
	ct := http.DetectContentType(head[:n])
	if !strings.HasPrefix(ct, "image/") {
		return "", ErrBadImage
	}
	_, _ = file.Seek(0, io.SeekStart)

	relDir := filepath.Join("products", productID)
	filename := uuid.NewString() + ext
	relPath := filepath.Join(relDir, filename)
	absDir := filepath.Join(s.uploadDir, relDir)
	absPath := filepath.Join(s.uploadDir, relPath)

	if err := os.MkdirAll(absDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir uploads dir: %w", err)
	}

	out, err := os.Create(absPath)
	if err != nil {
		return "", fmt.Errorf("create image file: %w", err)
	}
	defer func() {
		_ = out.Close()
		if err != nil {
			_ = os.Remove(absPath)
		}
	}()

	lr := io.LimitReader(file, MaxImageSize+1)
	written, err := io.Copy(out, lr)
	if err != nil {
		return "", fmt.Errorf("write image: %w", err)
	}
	if written > MaxImageSize {
		return "", ErrImageTooLarge
	}

	prev, hasPrev, err := s.repo.GetImagePath(ctx, productID)
	if err != nil {
		return "", fmt.Errorf("get previous image path: %w", err)
	}

	if err := s.repo.SetImagePath(ctx, productID, relPath); err != nil {
		return "", fmt.Errorf("update image path: %w", err)
	}

	if hasPrev && prev != "" && prev != relPath {
		_ = os.Remove(filepath.Join(s.uploadDir, prev))
	}

	return fmt.Sprintf("/api/products/%s/image", productID), nil
}

func (s *Service) Open(ctx context.Context, productID string) (f *os.File, filename string, mod time.Time, contentType string, etag string, err error) {
	rel, ok, err := s.repo.GetImagePath(ctx, productID)
	if err != nil {
		// If row doesn't exist, treat as product not found.
		return nil, "", time.Time{}, "", "", ErrProductNotFound
	}
	if !ok || strings.TrimSpace(rel) == "" {
		return nil, "", time.Time{}, "", "", ErrImageNotFound
	}

	if filepath.IsAbs(rel) {
		return nil, "", time.Time{}, "", "", ErrBadImage
	}

	abs := filepath.Clean(filepath.Join(s.uploadDir, rel))
	base := filepath.Clean(s.uploadDir) + string(os.PathSeparator)
	if abs != filepath.Clean(s.uploadDir) && !strings.HasPrefix(abs+string(os.PathSeparator), base) {
		return nil, "", time.Time{}, "", "", ErrBadImage
	}

	f, err = os.Open(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", time.Time{}, "", "", ErrImageNotFound
		}
		return nil, "", time.Time{}, "", "", fmt.Errorf("open image: %w", err)
	}

	st, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, "", time.Time{}, "", "", fmt.Errorf("stat image: %w", err)
	}
	if st.IsDir() {
		_ = f.Close()
		return nil, "", time.Time{}, "", "", ErrBadImage
	}

	contentType = mime.TypeByExtension(strings.ToLower(filepath.Ext(abs)))
	if contentType == "" {
		head := make([]byte, 512)
		n, _ := f.Read(head)
		contentType = http.DetectContentType(head[:n])
		_, _ = f.Seek(0, io.SeekStart)
	}

	mod = st.ModTime().UTC()
	etag = weakETag(st.Size(), mod, rel)

	return f, filepath.Base(abs), mod, contentType, etag, nil
}

func weakETag(size int64, mod time.Time, key string) string {
	// Include rel path in hash to avoid collisions across files with same size+mtime.
	h := sha256.Sum256([]byte(fmt.Sprintf("%d|%d|%s", size, mod.UnixNano(), key)))
	return "W/\"" + hex.EncodeToString(h[:8]) + "\""
}
