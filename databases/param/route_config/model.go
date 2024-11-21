package routeconfig

import "time"

type RouteConfig struct {
	ID        int64     `json:"id"`
	EndPoint  string    `json:"endpoint"`
	Url       string    `json:"url"`
	Status    int64     `json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (RouteConfig) TableName() string {
	return "route_config"
}