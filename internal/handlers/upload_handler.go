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

func (h *UploadHandler) UploadSoundFile(c *gin.Context) {
	ctx, span := h.logger.GetTracer().Start(c.Request.Context(), "UploadHandler.UploadSoundFile")
	defer span.End()

	wd, _ := os.Getwd()
	h.logger.Info("Current working directory: " + wd).WithTrace(ctx)

	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get file from form", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is requered"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		h.logger.Error("sound name is requered", err).WithTrace(ctx)
		c.JSON(http.StatusBadRequest, gin.H{"error": "sound name is requered"})
		return
	}

	uploadPath := "static/uploads"

	fullUploadPath := filepath.Join(wd, uploadPath)
	h.logger.Info("Full upload path: " + fullUploadPath).WithTrace(ctx)

	if err := ensureUploadDir(uploadPath); err != nil {
		h.logger.Error("failed to create upload directory", err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", name, fileExt)
	filePath := filepath.Join(uploadPath, fileName)

	h.logger.Info("Saving file to: " + filePath).WithTrace(ctx)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		h.logger.Error("failed to save file", err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		h.logger.Error("file was not created at: "+filePath, err).WithTrace(ctx)

		absPath, _ := filepath.Abs(filePath)
		h.logger.Info("Trying absolute path: " + absPath).WithTrace(ctx)
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			h.logger.Error("file was not created at absolute path either", err).WithTrace(ctx)
		} else {
			h.logger.Info("File found at absolute path: " + absPath).WithTrace(ctx)
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "file was not created"})
		return
	}

	fileInfo, _ := os.Stat(filePath)
	h.logger.Info("File successfully saved. Size: %d bytes", fileInfo.Size()).WithTrace(ctx)

	if err := h.service.UpdateSoundFile(ctx, name, fileName, filePath, file.Size); err != nil {
		h.logger.Error("failed to update sound file info", err).WithTrace(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update sound record"})
		return
	}

	h.logger.Info("file uploaded successfully", "filename", fileName, "size", file.Size, "path", filePath).WithTrace(ctx)
	c.JSON(http.StatusOK, gin.H{
		"message":   "file uploaded successfully",
		"filename":  fileName,
		"path":      filePath,
		"full_path": filepath.Join(wd, filePath),
	})
}

func ensureUploadDir(path string) error {
	return os.MkdirAll(path, 0775)
}
