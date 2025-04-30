package routes

import (
	"url-shortener/controllers"
	"url-shortener/metrics"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Aplica o middleware do Prometheus para todas as rotas
	router.Use(metrics.PrometheusMiddleware())

	// Endpoints da API
	urlController := controllers.UrlController{}
	router.POST("/shorten", urlController.CreateShortUrl)
	router.GET("/url/:shortUrl", urlController.GetOriginalUrl)

	// Endpoint para métricas do Prometheus
	router.GET("/metrics", metrics.PrometheusHandler())
}
