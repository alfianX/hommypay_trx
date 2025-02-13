package acqconfig

type AcqConfig struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (AcqConfig) TableName() string {
	return "acq_config"
}