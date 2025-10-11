package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"soundtube/internal/services"
	"soundtube/pkg"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service *services.SoundService
	logger  *pkg.CustomLogger
}

func NewUploadHandler(service *services.SoundService, logger *pkg.CustomLogger) *UploadHandler {
	return &UploadHandler{service: service, logger: logger}
}

// UploadSoundFile handles audio file upload
// @Summary Upload sound file
// @Description Upload an audio file for an existing sound record
// @Tags sounds
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Audio file to upload"
// @Param name formData string true "Sound name to associate with file"
// @Param request body UploadRequest true "Upload sound file"
// @Failure 400 {object} map[string]string "Missing file or sound name"
// @Failure 500 {object} map[string]string "File upload or database update failed"
// @Router /api/sounds/upload [post]
func (h *UploadHandler) UploadSoundFile(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "UploadHandler.UploadSoundFile")
	defer span.End()

	wd, _ := os.Getwd()
	projectRoot := filepath.Join(wd, "..", "..")
	h.logger.Info("Project root directory: " + projectRoot).WithTrace(ctx)

	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get file from form", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		h.logger.Error("sound name is required", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, gin.H{"error": "sound name is required"})
		return
	}

	uploadDir := "static/uploads"
	fullUploadDir := filepath.Join(projectRoot, uploadDir)
	h.logger.Info("Full upload directory: " + fullUploadDir).WithTrace(ctx)

	if err := ensureUploadDir(fullUploadDir); err != nil {
		h.logger.Error("failed to create upload directory", err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", name, fileExt)

	savePath := filepath.Join(fullUploadDir, fileName)

	dbFilePath := filepath.Join("uploads", fileName)

	h.logger.Info("Saving file to: " + savePath).WithTrace(ctx)
	h.logger.Info("DB file path: " + dbFilePath).WithTrace(ctx)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		h.logger.Error("failed to save file", err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		h.logger.Error("file was not created at: "+savePath, err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file was not created"})
		return
	}

	fileInfo, _ := os.Stat(savePath)
	h.logger.Info("File successfully saved. Size: %d bytes", fileInfo.Size()).WithTrace(ctx)

	if err := h.service.UpdateSoundFile(ctx, name, fileName, dbFilePath, file.Size); err != nil {
		h.logger.Error("failed to update sound file info", err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update sound record"})
		return
	}

	h.logger.Info("file uploaded successfully",
		"filename", fileName,
		"size", file.Size,
		"db_path", dbFilePath,
		"save_path", savePath).WithTrace(ctx)

	c.JSON(http.StatusOK, gin.H{
		"message":   "file uploaded successfully",
		"filename":  fileName,
		"path":      dbFilePath,
		"full_path": savePath,
	})
}

func ensureUploadDir(path string) error {
	return os.MkdirAll(path, 0775)
}
