package models

import (
	"time"

	"gorm.io/gorm"
)

// Url representa a estrutura de uma URL encurtada no banco de dados
type Url struct {
	ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
	ShortUrl    string         `json:"short_url" gorm:"unique;not null;index" example:"abc123"`
	OriginalUrl string         `json:"original_url" gorm:"not null" example:"https://www.exemplo.com.br/pagina-muito-grande-e-dificil-de-compartilhar"`
	CreatedAt   time.Time      `json:"created_at" example:"2025-05-01T12:34:56Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2025-05-01T12:34:56Z"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`                       // Ocultando DeletedAt do JSON para evitar problemas com Swagger
	Clicks      int            `json:"clicks" gorm:"default:0" example:"42"` // Contador de acessos à URL
}

// TableName especifica o nome da tabela no banco de dados
func (Url) TableName() string {
	return "urls"
}
