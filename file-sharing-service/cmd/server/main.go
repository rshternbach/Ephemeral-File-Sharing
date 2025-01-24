package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cyolo/file-sharing-service/config"
	"cyolo/file-sharing-service/internal/handlers"

	gorillaHandlers "github.com/gorilla/handlers"

	"cyolo/file-sharing-service/internal/utils"

	"github.com/gorilla/mux"
)

var logger = utils.Log

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/file", handlers.UploadFileHandler)
	r.HandleFunc("/v1/{file-url}", handlers.RetrieveFileHandler) // Updated route

	// Allow CORS
	corsHandler := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	go func() {
		for range time.Tick(config.CleanupInterval) {
			handlers.CleanupExpiredFiles()
		}
	}()

	if err := os.MkdirAll(config.UploadDir, os.ModePerm); err != nil {
		logger.Errorf("Failed to create upload directory: %v", err)
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	logger.Info("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(r)))
}
