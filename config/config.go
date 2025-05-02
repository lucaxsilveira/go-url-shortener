package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// DatabaseType representa o tipo de banco de dados
type DatabaseType string

const (
	PostgreSQL DatabaseType = "postgres"
	DynamoDB   DatabaseType = "dynamodb"
)

type Config struct {
	DatabaseURL  string
	DatabaseType DatabaseType
	Port         string
	AWSRegion    string
	AWSEndpoint  string // Para DynamoDB local
	RedisURL     string
}

// LoadEnv carrega as variáveis de ambiente do arquivo .env
func LoadEnv() {
	// Encontra o caminho do arquivo .env
	workDir, err := os.Getwd()
	if err != nil {
		log.Printf("Erro ao obter o diretório de trabalho: %v", err)
		return
	}

	// Carrega o arquivo .env
	envPath := filepath.Join(workDir, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Arquivo .env não encontrado ou não pode ser carregado: %v", err)
		log.Println("Usando variáveis de ambiente do sistema")
	} else {
		log.Println("Arquivo .env carregado com sucesso")
	}
}

// LoadConfig carrega as configurações da aplicação
func LoadConfig() *Config {
	// Carrega o arquivo .env
	LoadEnv()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/urlshortener"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Determina qual banco de dados usar
	dbType := os.Getenv("DATABASE_TYPE")
	var databaseType DatabaseType
	if dbType == "dynamodb" {
		databaseType = DynamoDB
	} else {
		databaseType = PostgreSQL
	}

	// Configurações AWS para DynamoDB
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	awsEndpoint := os.Getenv("AWS_ENDPOINT")

	// URL do Redis para cache
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	return &Config{
		DatabaseURL:  databaseURL,
		DatabaseType: databaseType,
		Port:         port,
		AWSRegion:    awsRegion,
		AWSEndpoint:  awsEndpoint,
		RedisURL:     redisURL,
	}
}
