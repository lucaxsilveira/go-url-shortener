package services

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ShortenerService struct {
	urlStore map[string]string
}

func NewShortenerService() *ShortenerService {
	return &ShortenerService{
		urlStore: make(map[string]string),
	}
}

func (s *ShortenerService) ShortenUrl(originalUrl string) string {
	shortUrl := generateShortUrl()
	fmt.Println("[ShortenUrl] Generated short URL:", shortUrl)
	fmt.Println("[ShortenUrl] Original URL:", originalUrl)

	// Inicializa o mapa se for nulo
	if s.urlStore == nil {
		s.urlStore = make(map[string]string)
	}

	s.urlStore[shortUrl] = originalUrl

	// print s.urlStore[shortUrl]
	fmt.Println("[ShortenUrl] URL Store:", s.urlStore)

	return shortUrl
}

func (s *ShortenerService) RetrieveUrl(shortUrl string) (string, bool) {
	// Verifica se o mapa é nulo antes de acessá-lo
	if s.urlStore == nil {
		return "", false
	}

	originalUrl, exists := s.urlStore[shortUrl]
	return originalUrl, exists
}

func generateShortUrl() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	rand.Seed(time.Now().UnixNano())
	var shortUrl strings.Builder
	for i := 0; i < length; i++ {
		shortUrl.WriteByte(charset[rand.Intn(len(charset))])
	}
	return shortUrl.String()
}
