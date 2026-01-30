package strutil

import "github.com/google/uuid"

func TruncateString(value string, maxSize int) string {
	if maxSize <= 0 {
		return ""
	}
	if len(value) <= maxSize {
		return value
	}
	const separator = "......"
	if maxSize <= len(separator) {
		return value[:maxSize]
	}
	headSize := (maxSize - len(separator)) / 2
	tailSize := maxSize - len(separator) - headSize
	return value[:headSize] + separator + value[len(value)-tailSize:]
}

func NewUUID() string {
	return uuid.NewString()
}
