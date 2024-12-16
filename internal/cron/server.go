package cron

import (
	"fmt"

	"github.com/alfianX/hommypay_trx/configs"
	"github.com/alfianX/hommypay_trx/databases"
	"github.com/alfianX/hommypay_trx/internal/cron/handlers"
)

type Server struct {
	cron handlers.Service
}

func NewServer() (*Server, error) {
	cnf, err := configs.NewParsedConfig()
	if err != nil {
		return nil, err
	}

	databaseTrx, err := databases.Connect(databases.ConfigDB{
		Host:     cnf.DatabaseTrx.Host,
		Port:     cnf.DatabaseTrx.Port,
		User:     cnf.DatabaseTrx.User,
		Password: cnf.DatabaseTrx.Password,
		Name:     cnf.DatabaseTrx.Name,
	})
	if err != nil {
		return nil, err
	}

	databaseParam, err := databases.Connect(databases.ConfigDB{
		Host:     cnf.DatabaseParam.Host,
		Port:     cnf.DatabaseParam.Port,
		User:     cnf.DatabaseParam.User,
		Password: cnf.DatabaseParam.Password,
		Name:     cnf.DatabaseParam.Name,
	})
	if err != nil {
		return nil, err
	}

	cron := handlers.NewHandler(databaseTrx, databaseParam)

	s := Server{
		cron: cron,
	}

	return &s, nil
}

func (s Server) Run() {
	fmt.Println("Cron job running...")
	s.cron.CronJob()
}