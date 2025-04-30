package config

import (
	"log"
	"url-shortener/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB inicializa a conexão com o banco de dados
func InitDB(config *Config) (*gorm.DB, error) {
	var err error

	DB, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
		return nil, err
	}

	// Migra os modelos para o banco de dados
	err = DB.AutoMigrate(&models.Url{})
	if err != nil {
		log.Fatalf("Erro ao migrar modelos: %v", err)
		return nil, err
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso")
	return DB, nil
}

// GetDB retorna a instância do banco de dados
func GetDB() *gorm.DB {
	return DB
}
