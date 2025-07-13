package models

import "time"

type URL struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URL       string    `gorm:"type:varchar(2048);not null" json:"url"`
	Status    string    `gorm:"type:enum('queued','running','done','error');default:'queued'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
