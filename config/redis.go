package config

import (
	"context"
	"time"
	"url-shortener/utils"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient      *redis.Client
	redisAvailable   bool
	redisInitialized bool
)

// GetRedisClient retorna a instância do cliente Redis
func GetRedisClient() *redis.Client {
	if !redisInitialized {
		// Carrega as configurações
		cfg := LoadConfig()

		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisURL,
			Password: "", // sem senha
			DB:       0,  // banco padrão
		})

		// Testa a conexão
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			utils.LogError(err, "Redis não está disponível. Algumas funcionalidades como cache e métricas detalhadas podem estar limitadas")
			redisAvailable = false
		} else {
			utils.InfoLogger.Println("Conexão com Redis estabelecida com sucesso")
			redisAvailable = true
		}

		redisInitialized = true
	}

	return redisClient
}

// IsRedisAvailable retorna true se o Redis estiver disponível
func IsRedisAvailable() bool {
	if !redisInitialized {
		GetRedisClient()
	}
	return redisAvailable
}
