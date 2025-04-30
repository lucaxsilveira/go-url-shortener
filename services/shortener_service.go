package services

import (
	"fmt"
	"url-shortener/config"
	"url-shortener/models"
	"url-shortener/utils"

	"gorm.io/gorm"
)

type ShortenerService struct {
	db *gorm.DB
}

func NewShortenerService() *ShortenerService {
	return &ShortenerService{
		db: config.GetDB(),
	}
}

func (s *ShortenerService) ListAllUrls() ([]models.Url, error) {
	var urls []models.Url
	result := s.db.Find(&urls)
	if result.Error != nil {
		utils.LogError(result.Error, "Erro ao buscar URLs do banco de dados")
		return nil, result.Error
	}

	utils.InfoLogger.Printf("Recuperadas %d URLs do banco de dados", len(urls))
	return urls, nil
}

func (s *ShortenerService) ShortenUrl(originalUrl string) string {
	// Verificar se a URL já existe no banco de dados
	var existingUrl models.Url
	result := s.db.Where("original_url = ?", originalUrl).First(&existingUrl)
	if result.RowsAffected > 0 {
		// Se a URL já existe, retorna o código curto existente
		utils.InfoLogger.Printf("URL já existe no banco de dados: %s", existingUrl.ShortUrl)
		return existingUrl.ShortUrl
	}

	// Gera um novo código curto
	shortUrl, err := utils.GenerateShortUrl()
	if err != nil {
		utils.LogError(err, "Erro ao gerar URL curta")
		return ""
	}

	utils.InfoLogger.Printf("Nova URL curta gerada: %s", shortUrl)

	// Cria um novo registro no banco de dados
	url := models.Url{
		ShortUrl:    shortUrl,
		OriginalUrl: originalUrl,
		Clicks:      0,
	}

	if err := s.db.Create(&url).Error; err != nil {
		utils.LogError(err, "Erro ao salvar URL no banco de dados")
		return ""
	}

	fmt.Println("[ShortenUrl] Generated short URL:", shortUrl)
	fmt.Println("[ShortenUrl] Original URL:", originalUrl)

	return shortUrl
}

func (s *ShortenerService) RetrieveUrl(shortUrl string) (string, bool) {
	var url models.Url

	result := s.db.Where("short_url = ?", shortUrl).First(&url)
	if result.Error != nil {
		utils.LogError(result.Error, "URL não encontrada no banco de dados")
		return "", false
	}

	// Incrementa o contador de cliques da URL
	s.db.Model(&url).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))

	return url.OriginalUrl, true
}
