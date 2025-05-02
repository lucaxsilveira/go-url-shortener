package repositories

import (
	"errors"
	"url-shortener/config"
)

// NewURLRepository cria um novo repositório de URLs com base no tipo de banco de dados configurado
func NewURLRepository() (URLRepository, error) {
	dbType := config.GetDBType()

	switch dbType {
	case config.PostgreSQL:
		return NewPostgresURLRepository(), nil
	case config.DynamoDB:
		return NewDynamoDBURLRepository(), nil
	default:
		return nil, errors.New("tipo de banco de dados não suportado")
	}
}
