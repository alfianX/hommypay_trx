package fdsconfig

import "time"

type FdsConfig struct {
	ID         int64     `json:"id"`
	FdsAddress string    `json:"fds_address"`
	CreatedAt  time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (FdsConfig) TableName() string {
	return "fds_config"
}