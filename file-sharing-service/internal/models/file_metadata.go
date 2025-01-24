package models

import "time"

type FileMetadata struct {
	Filename       string
	ExpirationTime time.Time
}
