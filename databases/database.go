package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/alfianX/hommypay_trx/pkg/round_robin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConfigDB struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func Connect(cfg ConfigDB) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	// sqlDB.SetMaxIdleConns(100)
	// sqlDB.SetMaxOpenConns(100)

	return db, nil
}

func InitIssuer(db *gorm.DB) (error) {
	var data []struct{
		ID 					int64 	`json:"id"`
		IssuerName			string	`json:"issuer_name"`
		IssuerConnType 		int64	`json:"issuer_conn_type"`
		IssuerService		string	`json:"issuer_service"` 
	}

	var id []int64
	var name []string
	var connType []int64
	var service []string

	result := db.WithContext(context.Background()).Table("issuer").Select("id", "issuer_name", "issuer_conn_type", "issuer_service").
				Where("issuer_type = ? AND status = ?", "DEBIT", 1).Find(&data)
	
	if result.Error != nil {
		return result.Error
	}

	for _, val := range data {
		id = append(id, val.ID)
		name = append(name, val.IssuerName)
		connType = append(connType, val.IssuerConnType)
		service = append(service, val.IssuerService)
	}

	round_robin.InitRoundRobin(id, name, connType, service)

	return nil
}