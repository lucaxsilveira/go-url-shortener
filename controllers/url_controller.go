package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"url-shortener/metrics"
	"url-shortener/models"
	"url-shortener/services"
	"url-shortener/utils"

	"github.com/gin-gonic/gin"
)

type UrlController struct {
	Service services.ShortenerService
}

func (uc *UrlController) CreateShortUrl(c *gin.Context) {
	var url models.Url

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

	if err := c.ShouldBindJSON(&url); err != nil {
		utils.LogError(err, "Erro ao fazer bind do JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Received URL:", url.OriginalUrl)

	shortUrl := uc.Service.ShortenUrl(url.OriginalUrl)
	utils.InfoLogger.Printf("URL encurtada gerada: %s -> %s", url.OriginalUrl, shortUrl)

	// Registra a métrica de URL encurtada
	metrics.RecordURLShortened()

	c.JSON(http.StatusCreated, gin.H{"short_url": shortUrl})
}

func (uc *UrlController) GetOriginalUrl(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	utils.LogRequest(c.Request.Method, c.Request.URL.Path, "")

	originalUrl, exists := uc.Service.RetrieveUrl(shortUrl)
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
