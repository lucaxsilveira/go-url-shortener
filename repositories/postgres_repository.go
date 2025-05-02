package repositories

import (
	"errors"
	"url-shortener/config"
	"url-shortener/models"

	"gorm.io/gorm"
)

// PostgresURLRepository implementa a interface URLRepository para PostgreSQL
type PostgresURLRepository struct {
	db *gorm.DB
}

// NewPostgresURLRepository cria um novo repositório PostgreSQL
func NewPostgresURLRepository() *PostgresURLRepository {
	return &PostgresURLRepository{
		db: config.GetDB(),
	}
}

// Create cria uma nova URL no banco de dados
func (r *PostgresURLRepository) Create(url *models.Url) error {
	return r.db.Create(url).Error
}

// GetByShortURL obtém uma URL pelo seu código curto
func (r *PostgresURLRepository) GetByShortURL(shortURL string) (*models.Url, error) {
	var url models.Url
	result := r.db.Where("short_url = ?", shortURL).First(&url)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("URL não encontrada")
		}
		return nil, result.Error
	}
	return &url, nil
}

// IncrementClicks incrementa o contador de cliques de uma URL
func (r *PostgresURLRepository) IncrementClicks(shortURL string) error {
	result := r.db.Model(&models.Url{}).
		Where("short_url = ?", shortURL).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("URL não encontrada")
	}
	return nil
}

// List retorna todas as URLs no banco de dados
func (r *PostgresURLRepository) List() ([]models.Url, error) {
	var urls []models.Url
	result := r.db.Find(&urls)
	return urls, result.Error
}

// Delete exclui uma URL pelo seu código curto
func (r *PostgresURLRepository) Delete(shortURL string) error {
	result := r.db.Where("short_url = ?", shortURL).Delete(&models.Url{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("URL não encontrada")
	}
	return nil
}
