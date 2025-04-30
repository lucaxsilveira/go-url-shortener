package main

import (
	"url-shortener/routes"
	"url-shortener/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializa o logger
	utils.InitLogger()

	router := gin.Default()

	routes.SetupRoutes(router)

	utils.InfoLogger.Println("Servidor inicializado na porta :8080")
	router.Run(":8080")
}
