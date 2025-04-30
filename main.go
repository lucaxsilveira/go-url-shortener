package main

import (
	"time"
	"url-shortener/cache"
	"url-shortener/config"
	"url-shortener/routes"
	"url-shortener/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializa o logger
	utils.InitLogger()

	// Carrega as configurações
	cfg := config.LoadConfig()

	// Inicializa o banco de dados
	_, err := config.InitDB(cfg)
	if err != nil {
		utils.ErrorLogger.Fatalf("Falha ao inicializar o banco de dados: %v", err)
	}

	// Inicializa o monitor do Redis (verifica métricas a cada 15 segundos)
	cache.StartRedisMonitor(15 * time.Second)

	router := gin.Default()

	routes.SetupRoutes(router)

	utils.InfoLogger.Println("Servidor inicializado na porta :" + cfg.Port)
	router.Run(":" + cfg.Port)
}
