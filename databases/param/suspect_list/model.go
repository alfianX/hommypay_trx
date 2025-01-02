package suspectlist

type SuspectList struct {
	ID     int64  `json:"id"`
	Mid    string `json:"mid"`
	Tid    string `json:"tid"`
	Trace  string `json:"trace"`
	Pan    string `json:"pan"`
	Date   string `json:"date"`
	Status string `json:"status"`
	Data   string `json:"data"`
}

func (SuspectList) TableName() string {
	return "suspect_list"
}