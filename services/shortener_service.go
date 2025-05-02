package services

import (
	"fmt"
	"url-shortener/models"
	"url-shortener/repositories"
	"url-shortener/utils"
)

type ShortenerService struct {
	repository repositories.URLRepository
}

func NewShortenerService() (*ShortenerService, error) {
	repository, err := repositories.NewURLRepository()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar repositório: %w", err)
	}

	return &ShortenerService{
		repository: repository,
	}, nil
}

func (s *ShortenerService) ListAllUrls() ([]models.Url, error) {
	// Iniciar medição do tempo total
	stopTotal := utils.MeasureExecutionTime("ListAllUrls (total)")
	defer func() {
		utils.InfoLogger.Printf("Tempo total de ListAllUrls: %v", stopTotal())
	}()

	urls, err := s.repository.List()
	if err != nil {
		utils.LogError(err, "Erro ao buscar URLs do banco de dados")
		return nil, err
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
	// Para simplificar, usamos uma abordagem que não é dependente de banco de dados específico
	// Recuperamos todas as URLs e procuramos por correspondência
	stopCheck := utils.MeasureExecutionTime("URL existence check")

	urls, err := s.repository.List()
	var existingUrl *models.Url

	if err == nil {
		for _, url := range urls {
			if url.OriginalUrl == originalUrl {
				existingUrl = &url
				break
			}
		}
	}

	checkTime := stopCheck()

	if existingUrl != nil {
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

	if err := s.repository.Create(&url); err != nil {
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
	// Iniciar medição do tempo total
	stopTotal := utils.MeasureExecutionTime("RetrieveUrl (total)")
	defer func() {
		utils.InfoLogger.Printf("Tempo total de RetrieveUrl: %v", stopTotal())
	}()

	// Iniciar medição do tempo de consulta ao banco de dados
	stopDB := utils.MeasureExecutionTime("Database lookup")

	url, err := s.repository.GetByShortURL(shortUrl)

	dbTime := stopDB()

	if err != nil {
		utils.LogError(err, fmt.Sprintf("URL não encontrada no banco de dados (em %v)", dbTime))
		return "", false
	}

	// Iniciar medição do tempo de atualização do contador de cliques
	stopUpdate := utils.MeasureExecutionTime("Click counter update")

	// Incrementa o contador de cliques da URL
	err = s.repository.IncrementClicks(shortUrl)
	if err != nil {
		utils.LogError(err, "Erro ao incrementar contador de cliques")
	}

	updateTime := stopUpdate()

	utils.InfoLogger.Printf("URL recuperada do banco de dados em %v, contador atualizado em %v", dbTime, updateTime)

	return url.OriginalUrl, true
}
