package handlers

import (
	"os"
	"path/filepath"
	"time"

	"cyolo/file-sharing-service/config"

	"cyolo/file-sharing-service/internal/utils"
)

var logger = utils.Log

func CleanupExpiredFiles() {
	metadataMutex.Lock()
	defer metadataMutex.Unlock()

	now := time.Now()
	logger.Infof("Running cleanup at: %v", now)
	for filePath, meta := range filesMetadata {
		logger.Infof("Checking file: %s, expiration time: %v", filePath, meta.ExpirationTime)
		if now.After(meta.ExpirationTime) {
			fullFilePath := filepath.Join(config.UploadDir, filePath)
			if err := os.Remove(fullFilePath); err == nil {
				delete(filesMetadata, filePath)
				logger.Infof("Deleted expired file: %s", fullFilePath)
			} else {
				logger.Errorf("Failed to delete file: %s, error: %v", fullFilePath, err)
			}
		}
	}
}
