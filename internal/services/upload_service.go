package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"portfolio-server/internal/config"
	"portfolio-server/internal/database"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type UploadService struct {
	minioClient *minio.Client
	bucket      string
	endpoint    string
	useSSL      bool
}

func NewUploadService() *UploadService {
	cfg := config.LoadConfig()
	return &UploadService{
		minioClient: database.GetMinioClient(),
		bucket:      cfg.MinIO.Bucket,
		endpoint:    cfg.MinIO.Endpoint,
		useSSL:      cfg.MinIO.UseSSL,
	}
}

type UploadResponse struct {
	URL      string `json:"url"`
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
}

// UploadImage 이미지를 MinIO에 업로드하고 URL 반환
func (s *UploadService) UploadImage(ctx context.Context, file *multipart.FileHeader) (*UploadResponse, error) {
	// 파일 열기
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// 고유한 파일명 생성
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)

	// 이미지 폴더에 저장
	objectName := fmt.Sprintf("images/%s", fileName)

	// 콘텐츠 타입 설정
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// MinIO에 업로드
	_, err = s.minioClient.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	// URL 생성
	protocol := "http"
	if s.useSSL {
		protocol = "https"
	}
	url := fmt.Sprintf("%s://%s/%s/%s", protocol, s.endpoint, s.bucket, objectName)

	return &UploadResponse{
		URL:      url,
		FileName: fileName,
		Size:     file.Size,
	}, nil
}

// DeleteImage MinIO에서 이미지 삭제
func (s *UploadService) DeleteImage(ctx context.Context, fileName string) error {
	objectName := fmt.Sprintf("images/%s", fileName)
	
	err := s.minioClient.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete from MinIO: %w", err)
	}

	return nil
}
