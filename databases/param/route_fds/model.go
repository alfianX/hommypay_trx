package routefds

type RouteFds struct {
	ID         int64  `json:"id"`
	Url        string `json:"url"`
	Keterangan string `json:"keterangan"`
	Status     int64  `json:"status"`
	Result     int64  `json:"result"`
	Data       string `json:"data"`
}

func (RouteFds) TableName() string {
	return "route_fds"
}
