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
