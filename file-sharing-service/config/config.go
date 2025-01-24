package config

import (
	"time"
)

const (
	UploadDir       = "uploads"
	DefaultTTL      = 1 * time.Minute
	CleanupInterval = 1 * time.Minute
)
