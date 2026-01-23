package database

import (
	"log"
	"portfolio-server/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinIO(cfg *config.MinIOConfig) error {
	var err error
	
	// MinIO 클라이언트 초기화
	MinioClient, err = minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return err
	}

	log.Printf("MinIO initialized successfully. Endpoint: %s", cfg.Endpoint)
	return nil
}

func GetMinioClient() *minio.Client {
	return MinioClient
}
