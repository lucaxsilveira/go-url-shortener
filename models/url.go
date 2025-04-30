package models

type Url struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	ShortUrl    string `json:"short_url" gorm:"unique;not null"`
	OriginalUrl string `json:"original_url" gorm:"not null"`
}
