package main

import (
	"context"
	"log"

	"github.com/alfianX/hommypay_trx/internal/sale"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(ctx context.Context) error {
	server, err := sale.NewServer()
	if err != nil {
		return err
	}

	err = server.Run(ctx)
	return err
}