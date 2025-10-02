package databases

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/alfianX/hommypay_trx/pkg/round_robin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ConfigDB struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// Custom logger struct
type FileBasedLogger struct {
	logDir string
}

// Implement LogMode
func (f FileBasedLogger) LogMode(level logger.LogLevel) logger.Interface {
	return f
}

// Implement Info
func (f FileBasedLogger) Info(ctx context.Context, msg string, data ...interface{}) {}

// Implement Warn
func (f FileBasedLogger) Warn(ctx context.Context, msg string, data ...interface{}) {}

// Implement Error - Fix by adding context parameter
func (f FileBasedLogger) Error(ctx context.Context, msg string, data ...interface{}) {}

// Implement Trace - Fix by using correct signature
func (f FileBasedLogger) Trace(ctx context.Context, _ time.Time, fc func() (string, int64), _ error) {
	sql, _ := fc()

	if strings.HasPrefix(strings.ToUpper(sql), "INSERT") || strings.HasPrefix(strings.ToUpper(sql), "UPDATE") || strings.HasPrefix(strings.ToUpper(sql), "DELETE") {
		// Get the current Go file name
		_, file, _, _ := runtime.Caller(3) // Adjusted depth to ensure the correct file
		fileParts := strings.Split(file, "/")
		fileName := strings.TrimSuffix(fileParts[len(fileParts)-2], ".go")

		currentTime := time.Now()
		gmtFormat := "20060102"
		dateString := currentTime.Format(gmtFormat)

		// Create a log file based on the Go file name
		logFilePath := fmt.Sprintf("%s/%s_%s.log", f.logDir, fileName, dateString)
		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()

		// Write query to the file

		// log.SetOutput(logFile)
		// log.Println(sql)
		logger := log.New(logFile, "", 0)
		logger.Println(sql + ";")
	}
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

	logDir := os.Getenv("LOG_QUERY")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal(err)
	}

	fileLogger := FileBasedLogger{logDir: logDir}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: fileLogger,
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

func InitIssuer(db *gorm.DB) error {
	var data []struct {
		ID             int64  `json:"id"`
		IssuerName     string `json:"issuer_name"`
		IssuerConnType int64  `json:"issuer_conn_type"`
		IssuerService  string `json:"issuer_service"`
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
