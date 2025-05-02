package routes

import (
	"url-shortener/controllers"
	"url-shortener/metrics"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Importação de documentos do Swagger para registro
	_ "url-shortener/docs"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Aplica o middleware do Prometheus para todas as rotas
	router.Use(metrics.PrometheusMiddleware())

	// Endpoints da API
	urlController := controllers.UrlController{}
	router.POST("/shorten", urlController.CreateShortUrl)
	router.GET("/url/:shortUrl", urlController.GetOriginalUrl)
	router.GET("/urls", urlController.ListAllUrls) // Rota para listar todas as URLs

	// Endpoint para métricas do Prometheus
	router.GET("/metrics", metrics.PrometheusHandler())

	// Endpoints do Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
