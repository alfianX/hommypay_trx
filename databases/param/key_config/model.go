package keyconfig

import "time"

type KeyConfig struct {
	ID        int64     `json:"id"`
	KeyType   string    `json:"key_type"`
	Value     string    `json:"value"`
	CreatedAt time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (KeyConfig) TableName() string {
	return "key_config"
}