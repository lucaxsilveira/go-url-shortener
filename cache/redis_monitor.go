// filepath: /Users/lucas.silveira/Documents/Study/encurtador/url-shortener/cache/redis_monitor.go
package cache

import (
	"context"
	"strconv"
	"time"

	"url-shortener/config"
	"url-shortener/metrics"
	"url-shortener/utils"

	"github.com/redis/go-redis/v9"
)

// StartRedisMonitor inicia um monitoramento periódico do Redis
func StartRedisMonitor(interval time.Duration) {
	go func() {
		for {
			monitorRedisHealth()
			time.Sleep(interval)
		}
	}()
	utils.InfoLogger.Println("Monitor de métricas do Redis iniciado com intervalo de", interval)
}

// monitorRedisHealth verifica a saúde do Redis e atualiza as métricas
func monitorRedisHealth() {
	client := config.GetRedisClient()
	ctx := context.Background()

	// Verifica se o Redis está funcionando com PING
	startTime := time.Now()
	pong, err := client.Ping(ctx).Result()
	latency := time.Since(startTime).Seconds()

	// Atualiza métrica de latência
	metrics.SetRedisLatency(latency)

	// Atualiza métrica de disponibilidade
	if err != nil || pong != "PONG" {
		metrics.SetRedisUp(false)
		utils.LogError(err, "Redis não está respondendo corretamente")
		return
	}

	metrics.SetRedisUp(true)

	// Obtém estatísticas de uso de memória
	info, err := client.Info(ctx, "memory").Result()
	if err != nil {
		utils.LogError(err, "Erro ao obter informações de memória do Redis")
	} else {
		// Extrai uso de memória do retorno de INFO MEMORY
		// O formato é "used_memory:1234\r\n..."
		memoryStr := extractRedisInfoValue(info, "used_memory")
		if memoryBytes, err := strconv.ParseFloat(memoryStr, 64); err == nil {
			metrics.SetRedisMemoryUsage(memoryBytes)
		}
	}

	// Obtém contagem total de chaves
	dbSize, err := client.DBSize(ctx).Result()
	if err != nil {
		utils.LogError(err, "Erro ao obter tamanho do banco Redis")
	} else {
		metrics.SetRedisKeysTotal(float64(dbSize))
	}

	// Calcula TTL médio para chaves existentes (amostragem de até 100 chaves)
	calculateAverageTTL(ctx, client)
}

// extractRedisInfoValue extrai um valor específico da saída do comando INFO
func extractRedisInfoValue(info string, key string) string {
	keyWithColon := key + ":"
	startIdx := 0

	// Encontra o índice onde a chave começa
	for i := 0; i < len(info); i++ {
		if i+len(keyWithColon) <= len(info) && info[i:i+len(keyWithColon)] == keyWithColon {
			startIdx = i + len(keyWithColon)
			break
		}
	}

	if startIdx == 0 {
		return "0" // Não encontrou a chave
	}

	// Encontra o fim do valor (até o \r\n)
	endIdx := startIdx
	for endIdx < len(info) && info[endIdx] != '\r' && info[endIdx] != '\n' {
		endIdx++
	}

	return info[startIdx:endIdx]
}

// calculateAverageTTL calcula o TTL médio das chaves no Redis
func calculateAverageTTL(ctx context.Context, client *redis.Client) {
	// Obtém todas as chaves (limitado a 100 para não impactar performance)
	keys, err := client.Keys(ctx, "*").Result()
	if err != nil {
		utils.LogError(err, "Erro ao obter chaves do Redis")
		return
	}

	if len(keys) == 0 {
		metrics.SetRedisCacheTTLAvg(0)
		return
	}

	// Limita o número de chaves verificadas para não impactar performance
	maxKeys := 100
	if len(keys) > maxKeys {
		keys = keys[:maxKeys]
	}

	var totalTTL float64
	var validTTLCount int

	for _, key := range keys {
		ttl, err := client.TTL(ctx, key).Result()
		if err != nil {
			utils.LogError(err, "Erro ao obter TTL para chave "+key)
			continue
		}

		// Apenas considera TTLs válidos (positivos)
		if ttl.Seconds() > 0 {
			totalTTL += ttl.Seconds()
			validTTLCount++
		}
	}

	// Calcula média
	var avgTTL float64
	if validTTLCount > 0 {
		avgTTL = totalTTL / float64(validTTLCount)
	}

	metrics.SetRedisCacheTTLAvg(avgTTL)
}
