package issuer

import "time"

type Issuer struct {
	ID               	int64     `json:"id"`
	IssuerName       	string    `json:"issuer_name"`
	IssuerType       	string    `json:"issuer_type"`
	IssuerConnType	 	int64	  `json:"issuer_conn_type"`
	IssuerService		string    `json:"issuer_service"`
	IssuerHost			string	  `json:"issuer_host"`
	Status           	int64     `json:"status"`
	CreatedAt        	time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	CreatedBy        	int64     `json:"created_by"`
	UpdatedAt        	time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
	UpdatedBy        	int64     `json:"updated_by"`
}

func (Issuer) TableName() string {
	return "issuer"
}