package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"cyolo/file-sharing-service/config"
	"cyolo/file-sharing-service/internal/models"
	"cyolo/file-sharing-service/internal/utils"

	"github.com/gorilla/mux"
)

var (
	filesMetadata = make(map[string]models.FileMetadata)
	metadataMutex sync.Mutex
)

func init() {
	loadMetadataFromFile()

}

// Generate a unique file ID
// The probability of not having a unique ID after 10 attempts is extremely low
// due to the large ID space (2^128 possible IDs). The probability of a collision
// in a single attempt is n / 2^128, where n is the number of existing IDs.
// For 10 attempts, the probability of not finding a unique ID is approximately (n / 2^128)^10.
func generateUniqueFileID() (string, error) {
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		fileID := utils.GenerateFileID()

		metadataMutex.Lock()
		if _, exists := filesMetadata[fileID]; !exists {
			metadataMutex.Unlock()
			return fileID, nil
		}
		metadataMutex.Unlock()
	}
	return "", fmt.Errorf("failed to generate a unique file ID after %d attempts", maxAttempts)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	utils.Log.Info("UploadFileHandler")

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		handleError(w, "Failed to parse multipart form", err, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		handleError(w, "Failed to read file", err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	ttl := getTTL(r.Header.Get("Retention-Time"))

	fileID, err := generateUniqueFileID()
	if err != nil {
		handleError(w, "Failed to generate unique file ID", err, http.StatusInternalServerError)
		return
	}

	if err := saveFile(file, fileID); err != nil {
		handleError(w, "Failed to save file", err, http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(time.Duration(ttl) * time.Minute)
	utils.Log.Infof("File uploaded: %s, expiration time: %v", fileID, expirationTime)

	saveMetadata(fileID, header.Filename, expirationTime)

	writeJSONResponse(w, map[string]string{"url": fileID})
}

func RetrieveFileHandler(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file-url"]

	fileMetadata, err := getFileMetadata(fileID)
	if err != nil {
		handleError(w, "File not found", err, http.StatusNotFound)
		return
	}

	if time.Now().After(fileMetadata.ExpirationTime) {
		handleError(w, "File has expired", fmt.Errorf("file has expired"), http.StatusNotFound)
		return
	}

	file, err := openFile(fileID)
	if err != nil {
		handleError(w, "Failed to open file", err, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		handleError(w, "Failed to send file", err, http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, message string, err error, statusCode int) {
	utils.Log.Error(message, ": ", err)
	http.Error(w, message, statusCode)
}

func getTTL(ttlStr string) int {
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil || ttl <= 0 {
		return int(config.DefaultTTL.Minutes())
	}
	return ttl
}

func saveFile(file io.Reader, fileID string) error {
	filePath := filepath.Join(config.UploadDir, fileID)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return err
	}
	return nil
}

func saveMetadata(fileID, filename string, expirationTime time.Time) {
	metadataMutex.Lock()
	defer metadataMutex.Unlock()

	fileMetadata := models.FileMetadata{
		Filename:       filename,
		ExpirationTime: expirationTime,
	}
	filesMetadata[fileID] = fileMetadata

	// Save metadata to file asynchronously
	go saveMetadataToFile(fileID, fileMetadata)
}

func writeJSONResponse(w http.ResponseWriter, response map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		handleError(w, "Failed to marshal JSON response", err, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func getFileMetadata(fileID string) (models.FileMetadata, error) {
	metadataMutex.Lock()
	defer metadataMutex.Unlock()
	fileMetadata, exists := filesMetadata[fileID]
	if !exists {
		return models.FileMetadata{}, fmt.Errorf("file metadata not found")
	}
	return fileMetadata, nil
}

func openFile(fileID string) (*os.File, error) {
	filePath := filepath.Join(config.UploadDir, fileID)
	return os.Open(filePath)
}
