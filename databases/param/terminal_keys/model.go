package terminalkeys

import "time"

type TerminalKeys struct {
	ID        int64     `json:"id"`
	Tid       string    `json:"tid"`
	KeyType   string    `json:"key_type"`
	Value     string     `json:"value"`
	CreatedAt time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}