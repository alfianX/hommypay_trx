package main

import (
	"log"

	"github.com/alfianX/hommypay_trx/internal/cron"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	server, err := cron.NewServer()
	if err != nil {
		return err
	}

	server.Run()

	return nil
}