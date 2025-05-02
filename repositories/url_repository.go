package repositories

import (
	"url-shortener/models"
)

// URLRepository define a interface para operações de repositório de URLs
type URLRepository interface {
	// Create cria uma nova URL no banco de dados
	Create(url *models.Url) error

	// GetByShortURL obtém uma URL pelo seu código curto
	GetByShortURL(shortURL string) (*models.Url, error)

	// IncrementClicks incrementa o contador de cliques de uma URL
	IncrementClicks(shortURL string) error

	// List retorna todas as URLs no banco de dados
	List() ([]models.Url, error)

	// Delete exclui uma URL pelo seu código curto
	Delete(shortURL string) error
}
