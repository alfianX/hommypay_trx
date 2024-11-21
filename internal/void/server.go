package void

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alfianX/hommypay_trx/configs"
	"github.com/alfianX/hommypay_trx/databases"
	"github.com/alfianX/hommypay_trx/internal"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logger *logrus.Logger
	router *gin.Engine
	config configs.Config
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

	databaseMerchant, err := databases.Connect(databases.ConfigDB{
		Host:     cnf.DatabaseMerchant.Host,
		Port:     cnf.DatabaseMerchant.Port,
		User:     cnf.DatabaseMerchant.User,
		Password: cnf.DatabaseMerchant.Password,
		Name:     cnf.DatabaseMerchant.Name,
	})
	if err != nil {
		return nil, err
	}

	log := internal.NewLogger()

	gin.SetMode(cnf.Mode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	RegisterRoutes(router, log, cnf, databaseTrx, databaseParam, databaseMerchant)

	s := Server{
		logger: log,
		config: cnf,
		router: router,
	}

	return &s, nil
}

func (s *Server) Run(ctx context.Context) error {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.ServerPort),
		Handler: s.router,
	}

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(stopServer)

	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		s.logger.Printf("REST API listening on PORT %d", s.config.ServerPort)
		serverErrors <- server.ListenAndServe()
	}(&wg)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting REST API server: %w", err)
	case <-stopServer:
		s.logger.Warn("server received STOP signal")

		err := server.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("graceful shutdown did not complete: %w", err)
		}
		wg.Wait()
		s.logger.Info("Server was shutdown gracefully")
	}
	return nil
}