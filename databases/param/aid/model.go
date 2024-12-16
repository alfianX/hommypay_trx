package aid

type AidList struct {
	ID      int64  `json:"id"`
	AppName string `json:"app_name"`
	Aid     string `json:"aid"`
}

func (AidList) TableName() string {
	return "aid_list"
}