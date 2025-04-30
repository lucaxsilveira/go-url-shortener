package controllers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"url-shortener/cache"
	"url-shortener/metrics"
	"url-shortener/models"
	"url-shortener/services"
	"url-shortener/utils"

	"github.com/gin-gonic/gin"
)

type UrlController struct{}

type CreateUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
}

func (uc *UrlController) CreateShortUrl(c *gin.Context) {
	var request CreateUrlRequest

	// Log informações da requisição
	c.Request.ParseForm()
	utils.LogPostData(c.Request.Method, c.Request.URL.Path, c.Request.PostForm)

	// Para logging completo do body JSON
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	if len(bodyBytes) > 0 {
		utils.LogRequest(c.Request.Method, c.Request.URL.Path, string(bodyBytes))
		// Restaurar o body para que o bind funcione
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.LogError(err, "Erro ao fazer bind do JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar se a URL é válida
	if !utils.IsValidUrl(request.OriginalUrl) {
		utils.LogError(nil, "URL inválida: "+request.OriginalUrl)
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL inválida"})
		return
	}

	// Criar o serviço
	service := services.NewShortenerService()

	// Encurtar a URL
	shortUrl := service.ShortenUrl(request.OriginalUrl)
	if shortUrl == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao encurtar URL"})
		return
	}

	utils.InfoLogger.Printf("URL encurtada gerada: %s -> %s", request.OriginalUrl, shortUrl)

	// Registra a métrica de URL encurtada
	metrics.RecordURLShortened()

	// Constrói a URL completa com o domínio do servidor
	baseURL := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	fullShortURL := scheme + "://" + baseURL + "/url/" + shortUrl

	c.JSON(http.StatusCreated, gin.H{
		"short_url":    shortUrl,
		"full_url":     fullShortURL,
		"original_url": request.OriginalUrl,
	})
}

func (uc *UrlController) ListAllUrls(c *gin.Context) {
	utils.LogRequest(c.Request.Method, c.Request.URL.Path, "")

	// Inicia a medição do tempo total
	stopTotal := utils.MeasureExecutionTime("ListAllUrls (total)")
	defer func() {
		duration := stopTotal()
		// Adiciona o tempo de execução ao header da resposta
		c.Header("X-Execution-Time", duration.String())
	}()

	// Chave única para o cache
	cacheKey := "list:all:urls"

	// Estrutura para armazenar o resultado
	var result struct {
		Urls      []models.Url `json:"urls"`
		FromCache bool         `json:"from_cache"`
	}

	// Tenta recuperar do cache primeiro
	stopCache := utils.MeasureExecutionTime("Cache retrieval")
	found, err := cache.Get(cacheKey, &result)
	cacheTime := stopCache()

	if err != nil {
		utils.LogError(err, "Erro ao verificar cache")
		// Continua com a execução normal em caso de erro no cache
	}

	if found {
		utils.InfoLogger.Printf("Dados recuperados do cache em %v", cacheTime)
		result.FromCache = true
		c.JSON(http.StatusOK, result)
		return
	}

	// Se não encontrou no cache, busca do banco de dados
	stopDB := utils.MeasureExecutionTime("Database query")

	// Criar o serviço
	service := services.NewShortenerService()

	// Recuperar todas as URLs
	urls, err := service.ListAllUrls()
	dbTime := stopDB()

	if err != nil {
		utils.LogError(err, "Erro ao listar URLs")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":          "Erro ao listar URLs",
			"execution_time": dbTime.String(),
		})
		return
	}

	// Prepara o resultado
	result.Urls = urls
	result.FromCache = false

	// Armazena no cache com TTL de 5 segundos
	go func() {
		if err := cache.Set(cacheKey, result, cache.DefaultTTL); err != nil {
			utils.LogError(err, "Erro ao armazenar no cache")
		}
	}()

	utils.InfoLogger.Printf("Retornando lista com %d URLs (do banco em %v)", len(urls), dbTime)
	c.JSON(http.StatusOK, result)
}

func (uc *UrlController) GetOriginalUrl(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	utils.LogRequest(c.Request.Method, c.Request.URL.Path, "")

	// Criar o serviço
	service := services.NewShortenerService()

	originalUrl, exists := service.RetrieveUrl(shortUrl)
	if !exists {
		utils.LogError(nil, "URL não encontrada: "+shortUrl)

		// Registra a métrica de URL não encontrada
		metrics.RecordURLNotFound()

		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Registra a métrica de acesso à URL
	metrics.RecordURLAccess()

	utils.InfoLogger.Printf("URL encontrada: %s -> %s", shortUrl, originalUrl)
	c.JSON(http.StatusOK, gin.H{"original_url": originalUrl})
}
