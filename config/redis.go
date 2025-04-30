package config

import (
	"context"
	"os"
	"time"
	"url-shortener/utils"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// GetRedisClient retorna a instância do cliente Redis
func GetRedisClient() *redis.Client {
	if redisClient == nil {
		// Inicializa o cliente Redis
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			redisURL = "localhost:6379"
		}

		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisURL,
			Password: "", // sem senha
			DB:       0,  // banco padrão
		})

		// Testa a conexão
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			utils.LogError(err, "Falha ao conectar ao Redis")
		} else {
			utils.InfoLogger.Println("Conexão com Redis estabelecida com sucesso")
		}
	}

	return redisClient
}
