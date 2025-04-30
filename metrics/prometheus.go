package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestCounter conta o número de requisições HTTP por método e código de status
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_shortener_http_requests_total",
			Help: "Total number of HTTP requests by method and status code",
		},
		[]string{"method", "status", "path"},
	)

	// RequestDuration mede a duração das requisições HTTP
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "url_shortener_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// URLShorteningCounter conta o número de URLs encurtadas
	URLShorteningCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "url_shortener_urls_shortened_total",
			Help: "Total number of URLs shortened",
		},
	)

	// URLAccessCounter conta o número de acessos a URLs encurtadas
	URLAccessCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "url_shortener_urls_access_total",
			Help: "Total number of short URL accesses",
		},
	)

	// URLNotFoundCounter conta o número de tentativas de acesso a URLs inexistentes
	URLNotFoundCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "url_shortener_urls_not_found_total",
			Help: "Total number of requests to non-existent short URLs",
		},
	)

	// CacheHitCounter conta o número de acertos no cache (cache hits)
	CacheHitCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_shortener_cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"operation", "key_prefix"},
	)

	// CacheMissCounter conta o número de erros no cache (cache misses)
	CacheMissCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_shortener_cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"operation", "key_prefix"},
	)

	// CacheOperationDuration mede a duração das operações de cache
	CacheOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "url_shortener_cache_operation_duration_seconds",
			Help:    "Cache operation duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"operation", "key_prefix"},
	)

	// RedisCacheTotalOperations conta o número total de operações realizadas no Redis
	RedisCacheTotalOperations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "url_shortener_redis_operations_total",
			Help: "Total number of operations performed on Redis",
		},
	)

	// RedisConnectionErrors conta o número de erros de conexão com o Redis
	RedisConnectionErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "url_shortener_redis_connection_errors_total",
			Help: "Total number of Redis connection errors",
		},
	)

	// RedisUp indica se o Redis está disponível (1) ou não (0)
	RedisUp = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "url_shortener_redis_up",
			Help: "Indicates whether Redis is up (1) or down (0)",
		},
	)

	// RedisLatency mede a latência da conexão com o Redis
	RedisLatency = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "url_shortener_redis_latency_seconds",
			Help: "Latency of Redis connection in seconds",
		},
	)

	// RedisMemoryUsage indica o uso de memória do Redis em bytes
	RedisMemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "url_shortener_redis_memory_usage_bytes",
			Help: "Redis memory usage in bytes",
		},
	)

	// RedisKeysTotal indica o número total de chaves no Redis
	RedisKeysTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "url_shortener_redis_keys_total",
			Help: "Total number of keys in Redis",
		},
	)

	// RedisCacheTTLAvg indica o TTL médio das chaves no Redis (em segundos)
	RedisCacheTTLAvg = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "url_shortener_redis_ttl_avg_seconds",
			Help: "Average TTL of keys in Redis in seconds",
		},
	)
)

// PrometheusMiddleware é um middleware do Gin para coletar métricas do Prometheus
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		// Se o path estiver vazio, use o path da requisição
		if path == "" {
			path = c.Request.URL.Path
		}

		// Prossegue com o processamento da requisição
		c.Next()

		// Registra a duração da requisição
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		// Incrementa o contador de requisições
		RequestCounter.WithLabelValues(method,
			string(rune(status)),
			path).Inc()

		// Observa a duração da requisição
		RequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	}
}

// PrometheusHandler retorna um handler para o endpoint /metrics
func PrometheusHandler() gin.HandlerFunc {
	handler := promhttp.Handler()

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// RecordURLShortened incrementa o contador de URLs encurtadas
func RecordURLShortened() {
	URLShorteningCounter.Inc()
}

// RecordURLAccess incrementa o contador de acessos a URLs encurtadas
func RecordURLAccess() {
	URLAccessCounter.Inc()
}

// RecordURLNotFound incrementa o contador de URLs não encontradas
func RecordURLNotFound() {
	URLNotFoundCounter.Inc()
}

// RecordCacheHit incrementa o contador de cache hits
func RecordCacheHit(operation string, keyPrefix string) {
	CacheHitCounter.WithLabelValues(operation, keyPrefix).Inc()
}

// RecordCacheMiss incrementa o contador de cache misses
func RecordCacheMiss(operation string, keyPrefix string) {
	CacheMissCounter.WithLabelValues(operation, keyPrefix).Inc()
}

// ObserveCacheOperationDuration registra a duração de uma operação de cache
func ObserveCacheOperationDuration(operation string, keyPrefix string, duration time.Duration) {
	CacheOperationDuration.WithLabelValues(operation, keyPrefix).Observe(duration.Seconds())
}

// RecordRedisCacheOperation incrementa o contador total de operações do Redis
func RecordRedisCacheOperation() {
	RedisCacheTotalOperations.Inc()
}

// RecordRedisConnectionError incrementa o contador de erros de conexão com o Redis
func RecordRedisConnectionError() {
	RedisConnectionErrors.Inc()
}

// SetRedisUp atualiza o status de disponibilidade do Redis
func SetRedisUp(isUp bool) {
	if isUp {
		RedisUp.Set(1)
	} else {
		RedisUp.Set(0)
	}
}

// SetRedisLatency atualiza a latência da conexão com o Redis
func SetRedisLatency(latency float64) {
	RedisLatency.Set(latency)
}

// SetRedisMemoryUsage atualiza o uso de memória do Redis
func SetRedisMemoryUsage(bytes float64) {
	RedisMemoryUsage.Set(bytes)
}

// SetRedisKeysTotal atualiza o número total de chaves no Redis
func SetRedisKeysTotal(count float64) {
	RedisKeysTotal.Set(count)
}

// SetRedisCacheTTLAvg atualiza o TTL médio das chaves no Redis
func SetRedisCacheTTLAvg(ttl float64) {
	RedisCacheTTLAvg.Set(ttl)
}
