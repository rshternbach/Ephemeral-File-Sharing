package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cyolo/file-sharing-service/config"
	"cyolo/file-sharing-service/internal/models"
	"cyolo/file-sharing-service/internal/utils"
)

var (
	metadataFile = filepath.Join(config.UploadDir, "metadata.txt")
	fileMutex    sync.Mutex
)

func saveMetadataToFile(fileID string, metadata models.FileMetadata) {
	utils.Log.Infof("Saving metadata for file ID: %s", fileID)
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(metadataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.Log.Error("Failed to open metadata file for writing: ", err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(metadata)
	if err != nil {
		utils.Log.Error("Failed to marshal metadata for file ID: ", fileID, ", error: ", err)
		return
	}

	if _, err := file.WriteString(fmt.Sprintf("%s %s\n", fileID, data)); err != nil {
		utils.Log.Error("Failed to write metadata to file for file ID: ", fileID, ", error: ", err)
	} else {
		utils.Log.Infof("Successfully saved metadata for file ID: %s", fileID)
	}
}

func loadMetadataFromFile() {
	utils.Log.Info("Loading metadata from file")
	file, err := os.Open(metadataFile)
	if err != nil {
		if os.IsNotExist(err) {
			utils.Log.Info("Metadata file does not exist, starting with empty metadata")
			return
		}
		utils.Log.Error("Failed to open metadata file for reading: ", err)
		return
	}
	defer file.Close()

	lines, err := readLastNLines(file, 1000)
	if err != nil {
		utils.Log.Error("Failed to read last 1000 lines from metadata file: ", err)
		return
	}

	processMetadataLines(lines)
	utils.Log.Infof("Finished loading metadata from file, loaded %d records", len(lines))
}

func readLastNLines(file *os.File, n int) ([]string, error) {
	utils.Log.Infof("Reading last %d lines from metadata file", n)
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(lines) == n {
			lines = lines[1:]
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	utils.Log.Infof("Successfully read last %d lines from metadata file", len(lines))
	return lines, nil
}

func processMetadataLines(lines []string) {
	utils.Log.Infof("Processing %d metadata lines", len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			utils.Log.Error("Invalid metadata line: ", line)
			continue
		}

		fileID := parts[0]
		var metadata models.FileMetadata
		if err := json.Unmarshal([]byte(parts[1]), &metadata); err != nil {
			utils.Log.Error("Failed to unmarshal metadata for file ID: ", fileID, ", error: ", err)
			continue
		}

		if fileExists(filepath.Join(config.UploadDir, fileID)) {
			updateInMemoryMetadata(fileID, metadata)
			utils.Log.Infof("Loaded metadata for file ID: %s", fileID)
		} else {
			utils.Log.Infof("File does not exist for file ID: %s, skipping", fileID)
		}
	}
	utils.Log.Info("Finished processing metadata lines")
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	exists := err == nil
	if exists {
		utils.Log.Infof("File exists: %s", filePath)
	} else {
		utils.Log.Infof("File does not exist: %s", filePath)
	}
	return exists
}

func updateInMemoryMetadata(fileID string, metadata models.FileMetadata) {
	utils.Log.Infof("Updating in-memory metadata for file ID: %s", fileID)
	metadataMutex.Lock()
	defer metadataMutex.Unlock()
	filesMetadata[fileID] = metadata
	utils.Log.Infof("Successfully updated in-memory metadata for file ID: %s", fileID)
}
