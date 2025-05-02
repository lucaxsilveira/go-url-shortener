package cache

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
	"url-shortener/config"
	"url-shortener/metrics"
	"url-shortener/utils"
)

const (
	// DefaultTTL é o tempo padrão de expiração do cache
	DefaultTTL = 5 * time.Second
)

var (
	// ErrRedisUnavailable é retornado quando o Redis não está disponível
	ErrRedisUnavailable = errors.New("redis não está disponível")
)

// getKeyPrefix extrai o prefixo da chave para fins de monitoramento
func getKeyPrefix(key string) string {
	parts := strings.Split(key, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

// Set armazena um valor no cache com TTL específico
func Set(key string, value interface{}, ttl time.Duration) error {
	// Verifica se o Redis está disponível
	if !config.IsRedisAvailable() {
		utils.InfoLogger.Printf("Cache desativado: Redis não está disponível")
		return ErrRedisUnavailable
	}

	startTime := time.Now()
	keyPrefix := getKeyPrefix(key)
	metrics.RecordRedisCacheOperation()

	ctx := context.Background()
	client := config.GetRedisClient()

	// Converte o valor para JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		utils.LogError(err, "Erro ao converter valor para JSON")
		return err
	}

	// Armazena no Redis
	err = client.Set(ctx, key, jsonValue, ttl).Err()
	if err != nil {
		utils.LogError(err, "Erro ao armazenar no cache")
		metrics.RecordRedisConnectionError()
		return err
	}

	// Registra a duração da operação
	duration := time.Since(startTime)
	metrics.ObserveCacheOperationDuration("set", keyPrefix, duration)

	utils.InfoLogger.Printf("Valor armazenado em cache: %s (TTL: %v)", key, ttl)
	return nil
}

// Get recupera um valor do cache
func Get(key string, result interface{}) (bool, error) {
	// Verifica se o Redis está disponível
	if !config.IsRedisAvailable() {
		utils.InfoLogger.Printf("Cache desativado: Redis não está disponível")
		return false, ErrRedisUnavailable
	}

	startTime := time.Now()
	keyPrefix := getKeyPrefix(key)
	metrics.RecordRedisCacheOperation()

	ctx := context.Background()
	client := config.GetRedisClient()

	// Recupera do Redis
	jsonValue, err := client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			// Cache miss - chave não encontrada
			metrics.RecordCacheMiss("get", keyPrefix)
			metrics.ObserveCacheOperationDuration("get", keyPrefix, time.Since(startTime))
			return false, nil
		}
		utils.LogError(err, "Erro ao recuperar do cache")
		metrics.RecordRedisConnectionError()
		return false, err
	}

	// Converte o JSON para o tipo desejado
	err = json.Unmarshal([]byte(jsonValue), result)
	if err != nil {
		utils.LogError(err, "Erro ao converter JSON para o tipo desejado")
		return false, err
	}

	// Cache hit
	metrics.RecordCacheHit("get", keyPrefix)
	metrics.ObserveCacheOperationDuration("get", keyPrefix, time.Since(startTime))

	utils.InfoLogger.Printf("Valor recuperado do cache: %s", key)
	return true, nil
}

// Delete remove um valor do cache
func Delete(key string) error {
	// Verifica se o Redis está disponível
	if !config.IsRedisAvailable() {
		utils.InfoLogger.Printf("Cache desativado: Redis não está disponível")
		return ErrRedisUnavailable
	}

	startTime := time.Now()
	keyPrefix := getKeyPrefix(key)
	metrics.RecordRedisCacheOperation()

	ctx := context.Background()
	client := config.GetRedisClient()

	err := client.Del(ctx, key).Err()
	if err != nil {
		utils.LogError(err, "Erro ao remover do cache")
		metrics.RecordRedisConnectionError()
		return err
	}

	// Registra a duração da operação
	metrics.ObserveCacheOperationDuration("delete", keyPrefix, time.Since(startTime))

	utils.InfoLogger.Printf("Valor removido do cache: %s", key)
	return nil
}
