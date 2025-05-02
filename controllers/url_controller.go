package controllers

import (
	"bytes"
	"io"
	"net/http"
	"url-shortener/cache"
	"url-shortener/config"
	"url-shortener/metrics"
	"url-shortener/models"
	"url-shortener/services"
	"url-shortener/utils"

	"github.com/gin-gonic/gin"
)

type UrlController struct{}

type CreateUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required" example:"https://www.exemplo.com.br/pagina-muito-grande-e-dificil-de-compartilhar"`
}

type CreateUrlResponse struct {
	ShortUrl      string `json:"short_url" example:"abc123"`
	FullUrl       string `json:"full_url" example:"http://localhost:8080/url/abc123"`
	OriginalUrl   string `json:"original_url" example:"https://www.exemplo.com.br/pagina-muito-grande-e-dificil-de-compartilhar"`
	ExecutionTime string `json:"execution_time" example:"12.345ms"`
}

type ErrorResponse struct {
	Error         string `json:"error" example:"URL inválida"`
	ExecutionTime string `json:"execution_time,omitempty" example:"5.678ms"`
}

type GetOriginalUrlResponse struct {
	OriginalUrl   string `json:"original_url" example:"https://www.exemplo.com.br/pagina-muito-grande-e-dificil-de-compartilhar"`
	ExecutionTime string `json:"execution_time" example:"3.456ms"`
}

type ListUrlsResponse struct {
	Urls          []models.Url `json:"urls"`
	FromCache     bool         `json:"from_cache" example:"true"`
	ExecutionTime string       `json:"execution_time" example:"2.345ms"`
	CacheEnabled  bool         `json:"cache_enabled" example:"true"`
}

// CreateShortUrl godoc
// @Summary      Encurtar uma URL
// @Description  Cria uma versão curta de uma URL longa
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        request body CreateUrlRequest true "URL original para encurtar"
// @Success      201  {object}  CreateUrlResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /shorten [post]
func (uc *UrlController) CreateShortUrl(c *gin.Context) {
	var request CreateUrlRequest

	// Inicia a medição do tempo total
	stopTotal := utils.MeasureExecutionTime("CreateShortUrl (total)")
	defer func() {
		duration := stopTotal()
		// Adiciona o tempo de execução ao header da resposta
		c.Header("X-Execution-Time", duration.String())
	}()

	// Log informações da requisição
	c.Request.ParseForm()
	utils.LogPostData(c.Request.Method, c.Request.URL.Path, c.Request.PostForm)

	// Para logging completo do body JSON
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	if len(bodyBytes) > 0 {
		utils.LogRequest(c.Request.Method, c.Request.URL.Path, string(bodyBytes))
		// Restaurar o body para que o bind funcione
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
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
	service, err := services.NewShortenerService()
	if err != nil {
		utils.LogError(err, "Erro ao criar serviço de encurtamento")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	// Medir o tempo da operação de encurtamento
	stopShorten := utils.MeasureExecutionTime("URL shortening")

	// Encurtar a URL
	shortUrl := service.ShortenUrl(request.OriginalUrl)
	shortenTime := stopShorten()

	if shortUrl == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":          "Erro ao encurtar URL",
			"execution_time": shortenTime.String(),
		})
		return
	}

	utils.InfoLogger.Printf("URL encurtada gerada: %s -> %s (em %v)", request.OriginalUrl, shortUrl, shortenTime)

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
		"short_url":      shortUrl,
		"full_url":       fullShortURL,
		"original_url":   request.OriginalUrl,
		"execution_time": shortenTime.String(),
	})
}

// ListAllUrls godoc
// @Summary      Listar todas as URLs
// @Description  Retorna todas as URLs encurtadas no sistema
// @Tags         urls
// @Produce      json
// @Success      200  {object}  ListUrlsResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /urls [get]
func (uc *UrlController) ListAllUrls(c *gin.Context) {
	utils.LogRequest(c.Request.Method, c.Request.URL.Path, "")

	// Inicia a medição do tempo total
	stopTotal := utils.MeasureExecutionTime("ListAllUrls (total)")
	defer func() {
		duration := stopTotal()
		// Adiciona o tempo de execução ao header da resposta
		c.Header("X-Execution-Time", duration.String())
	}()

	// Estrutura para armazenar o resultado
	var result struct {
		Urls          []models.Url `json:"urls"`
		FromCache     bool         `json:"from_cache"`
		ExecutionTime string       `json:"execution_time"`
		CacheEnabled  bool         `json:"cache_enabled"`
	}

	// Define se o cache está ativado
	result.CacheEnabled = config.IsRedisAvailable()

	// Se o Redis estiver disponível, tenta recuperar do cache
	if result.CacheEnabled {
		// Chave única para o cache
		cacheKey := "list:all:urls"

		// Tenta recuperar do cache primeiro
		stopCache := utils.MeasureExecutionTime("Cache retrieval")
		found, err := cache.Get(cacheKey, &result)
		cacheTime := stopCache()

		if err != nil && err != cache.ErrRedisUnavailable {
			utils.LogError(err, "Erro ao verificar cache")
			// Continua com a execução normal em caso de erro no cache
		}

		if found {
			utils.InfoLogger.Printf("Dados recuperados do cache em %v", cacheTime)
			result.FromCache = true
			c.JSON(http.StatusOK, result)
			return
		}
	} else {
		utils.InfoLogger.Printf("Cache desativado: Redis não está disponível")
	}

	// Se não encontrou no cache ou o cache está desativado, busca do banco de dados
	stopDB := utils.MeasureExecutionTime("Database query")

	// Criar o serviço
	service, err := services.NewShortenerService()
	if err != nil {
		utils.LogError(err, "Erro ao criar serviço de encurtamento")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

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
	result.ExecutionTime = dbTime.String()

	// Armazena no cache com TTL de 5 segundos (apenas se o Redis estiver disponível)
	if result.CacheEnabled {
		go func() {
			err := cache.Set("list:all:urls", result, cache.DefaultTTL)
			if err != nil && err != cache.ErrRedisUnavailable {
				utils.LogError(err, "Erro ao armazenar no cache")
			}
		}()
	}

	utils.InfoLogger.Printf("Retornando lista com %d URLs (do banco em %v)", len(urls), dbTime)
	c.JSON(http.StatusOK, result)
}

// GetOriginalUrl godoc
// @Summary      Obter URL original
// @Description  Retorna a URL original a partir do código curto
// @Tags         urls
// @Produce      json
// @Param        shortUrl path string true "Código da URL curta" example:"abc123"
// @Success      200  {object}  GetOriginalUrlResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /url/{shortUrl} [get]
func (uc *UrlController) GetOriginalUrl(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	utils.LogRequest(c.Request.Method, c.Request.URL.Path, "")

	// Inicia a medição do tempo total
	stopTotal := utils.MeasureExecutionTime("GetOriginalUrl (total)")
	defer func() {
		duration := stopTotal()
		// Adiciona o tempo de execução ao header da resposta
		c.Header("X-Execution-Time", duration.String())
	}()

	// Criar o serviço
	service, err := services.NewShortenerService()
	if err != nil {
		utils.LogError(err, "Erro ao criar serviço de encurtamento")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	// Medir o tempo da recuperação da URL
	stopRetrieve := utils.MeasureExecutionTime("URL retrieval")
	originalUrl, exists := service.RetrieveUrl(shortUrl)
	retrieveTime := stopRetrieve()

	if !exists {
		utils.LogError(nil, "URL não encontrada: "+shortUrl)

		// Registra a métrica de URL não encontrada
		metrics.RecordURLNotFound()

		c.JSON(http.StatusNotFound, gin.H{
			"error":          "URL not found",
			"execution_time": retrieveTime.String(),
		})
		return
	}

	// Registra a métrica de acesso à URL
	metrics.RecordURLAccess()

	utils.InfoLogger.Printf("URL encontrada: %s -> %s (em %v)", shortUrl, originalUrl, retrieveTime)
	c.JSON(http.StatusOK, gin.H{
		"original_url":   originalUrl,
		"execution_time": retrieveTime.String(),
	})
}
