package handlers

import (
	"log"
	"net/http"
	"portfolio-server/internal/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService *services.UploadService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		uploadService: services.NewUploadService(),
	}
}

// UploadImage godoc
// @Summary 이미지 업로드
// @Description 이미지를 MinIO에 업로드하고 URL을 반환합니다
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "업로드할 이미지 파일"
// @Success 200 {object} services.UploadResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 파일 가져오기
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "이미지 파일을 찾을 수 없습니다",
		})
		return
	}

	// 파일 크기 제한 (10MB)
	maxSize := int64(10 << 20) // 10MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "파일 크기는 10MB를 초과할 수 없습니다",
		})
		return
	}

	// 이미지 파일 타입 검증
	contentType := file.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "허용되지 않는 파일 형식입니다. (jpeg, jpg, png, gif, webp만 가능)",
		})
		return
	}

	// 업로드 처리
	result, err := h.uploadService.UploadImage(c.Request.Context(), file)
	if err != nil {
		log.Printf("Upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "이미지 업로드에 실패했습니다",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// DeleteImage godoc
// @Summary 이미지 삭제
// @Description MinIO에서 이미지를 삭제합니다
// @Tags upload
// @Accept json
// @Produce json
// @Param fileName query string true "삭제할 파일명"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /upload/image [delete]
func (h *UploadHandler) DeleteImage(c *gin.Context) {
	fileName := c.Query("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "fileName이 필요합니다",
		})
		return
	}

	err := h.uploadService.DeleteImage(c.Request.Context(), fileName)
	if err != nil {
		log.Printf("Delete error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "이미지 삭제에 실패했습니다",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "이미지가 삭제되었습니다",
	})
}
