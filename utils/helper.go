package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/url"
	"time"
)

// GenerateShortUrl generates a unique short URL string.
func GenerateShortUrl() (string, error) {
	b := make([]byte, 6) // 6 bytes will give us a 8-character string
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// IsValidUrl checks if the provided string is a valid URL.
func IsValidUrl(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}

// MeasureExecutionTime mede o tempo de execução de uma função e retorna a duração
func MeasureExecutionTime(funcName string) func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		duration := time.Since(start)
		InfoLogger.Printf("Tempo de execução de %s: %v", funcName, duration)
		return duration
	}
}
