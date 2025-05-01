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
	// Iniciar medição do tempo total
	stopTotal := utils.MeasureExecutionTime("ShortenUrl (total)")
	defer func() {
		utils.InfoLogger.Printf("Tempo total de ShortenUrl: %v", stopTotal())
	}()

	// Verificar se a URL já existe no banco de dados
	stopCheck := utils.MeasureExecutionTime("URL existence check")
	var existingUrl models.Url
	result := s.db.Where("original_url = ?", originalUrl).First(&existingUrl)
	checkTime := stopCheck()

	if result.RowsAffected > 0 {
		// Se a URL já existe, retorna o código curto existente
		utils.InfoLogger.Printf("URL já existe no banco de dados: %s (verificado em %v)", existingUrl.ShortUrl, checkTime)
		return existingUrl.ShortUrl
	}

	// Gera um novo código curto
	stopGenerate := utils.MeasureExecutionTime("Short code generation")
	shortUrl, err := utils.GenerateShortUrl()
	generateTime := stopGenerate()

	if err != nil {
		utils.LogError(err, fmt.Sprintf("Erro ao gerar URL curta (em %v)", generateTime))
		return ""
	}

	utils.InfoLogger.Printf("Nova URL curta gerada: %s (em %v)", shortUrl, generateTime)

	// Cria um novo registro no banco de dados
	stopSave := utils.MeasureExecutionTime("Database save")
	url := models.Url{
		ShortUrl:    shortUrl,
		OriginalUrl: originalUrl,
		Clicks:      0,
	}

	if err := s.db.Create(&url).Error; err != nil {
		utils.LogError(err, fmt.Sprintf("Erro ao salvar URL no banco de dados (em %v)", stopSave()))
		return ""
	}
	saveTime := stopSave()

	utils.InfoLogger.Printf("URL salva no banco de dados em %v", saveTime)

	fmt.Println("[ShortenUrl] Generated short URL:", shortUrl)
	fmt.Println("[ShortenUrl] Original URL:", originalUrl)

	return shortUrl
}

func (s *ShortenerService) RetrieveUrl(shortUrl string) (string, bool) {
	var url models.Url

	// Iniciar medição do tempo de consulta ao banco de dados
	stopDB := utils.MeasureExecutionTime("Database lookup")

	result := s.db.Where("short_url = ?", shortUrl).First(&url)

	dbTime := stopDB()

	if result.Error != nil {
		utils.LogError(result.Error, fmt.Sprintf("URL não encontrada no banco de dados (em %v)", dbTime))
		return "", false
	}

	// Iniciar medição do tempo de atualização do contador de cliques
	stopUpdate := utils.MeasureExecutionTime("Click counter update")

	// Incrementa o contador de cliques da URL
	s.db.Model(&url).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))

	updateTime := stopUpdate()

	utils.InfoLogger.Printf("URL recuperada do banco de dados em %v, contador atualizado em %v", dbTime, updateTime)

	return url.OriginalUrl, true
}
