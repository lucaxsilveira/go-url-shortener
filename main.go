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

	// Log o tipo de banco de dados que será usado
	utils.InfoLogger.Printf("Utilizando banco de dados: %s", cfg.DatabaseType)

	// Inicializa o banco de dados
	_, _, err := config.InitDB(cfg)
	if err != nil {
		utils.ErrorLogger.Fatalf("Falha ao inicializar o banco de dados: %v", err)
	}

	if cfg.DatabaseType == config.PostgreSQL {
		utils.InfoLogger.Printf("PostgreSQL inicializado com sucesso")
	} else if cfg.DatabaseType == config.DynamoDB {
		utils.InfoLogger.Printf("DynamoDB inicializado com sucesso")
	}

	// Inicializa o monitor do Redis (verifica métricas a cada 15 segundos)
	cache.StartRedisMonitor(15 * time.Second)

	router := gin.Default()

	routes.SetupRoutes(router)

	utils.InfoLogger.Println("Servidor inicializado na porta :" + cfg.Port)
	router.Run(":" + cfg.Port)
}
