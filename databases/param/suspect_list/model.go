package suspectlist

type SuspectList struct {
	ID     int64  `json:"id"`
	Mid    string `json:"MID"`
	Tid    string `json:"TID"`
	Trace  string `json:"trace"`
	Pan    string `json:"PAN"`
	Date   string `json:"date"`
	Status string `json:"status"`
	Data   string `json:"data"`
}

func (SuspectList) TableName() string {
	return "suspect_list"
}