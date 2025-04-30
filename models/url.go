package models

import (
	"time"

	"gorm.io/gorm"
)

// Url representa a estrutura de uma URL encurtada no banco de dados
type Url struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ShortUrl    string         `json:"short_url" gorm:"unique;not null;index"`
	OriginalUrl string         `json:"original_url" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Clicks      int            `json:"clicks" gorm:"default:0"` // Contador de acessos à URL
}

// TableName especifica o nome da tabela no banco de dados
func (Url) TableName() string {
	return "urls"
}
