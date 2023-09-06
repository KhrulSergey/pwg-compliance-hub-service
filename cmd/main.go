package main

import (
	"github.com/valyala/fasthttp"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/config"
	grpc "gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/grpc/server"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/router"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/service"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/storage"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/logger"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	appConfig   *config.AppConfig
	dbConfig    *config.DBConfig
	appLogger   logger.Logger
	dbConnector *gorm.DB

	//Storage
	complianceRepository storage.ComplianceRepository

	//Service
	externalComplianceService    service.ExternalComplianceService
	accountOperatorServiceClient service.AccountOperatorServiceClient
	complianceService            service.ComplianceService

	//REST
	complianceController *rest.ComplianceController
)

type shutdown = func() error

func main() {
	log.Println("INFO: start initialization")

	if err := initializeDependencies(); err != nil {
		log.Fatalf("FATAL: initialization failed: %v", err)
	}

	log.Println("INFO: initialization finished")

	errCh := make(chan error, 1)
	complianceService.StartScheduler(errCh)

	grpcShutdownFunc, err := startGRPCServer(errCh)
	if err != nil {
		log.Fatalf("FATAL: failed to start GRPC server: %v", err)
	}
	httpServerShutdown, err := startHTTPServer(errCh)
	if err != nil {
		log.Fatalf("FATAL: failed to start HTTP server: %v", err)
	}

	log.Println("INFO: server ready to accept connections")

	waitServers(
		errCh,
		grpcShutdownFunc,
		httpServerShutdown,
	)
}

func initializeDependencies() error {
	var err error

	appLogger, err = logger.NewRelease()
	defer func() {
		_ = appLogger.Flush()
	}() // flushes buffer, if any

	// Define new appConfig
	appConfig, err = config.InitAppConfig()
	if err != nil {
		appLogger.Fatalf("Unable to read configuration %v", err)
	}

	//update logger level
	if appConfig.LoggerMode == "DEBUG" {
		appLogger, _ = logger.NewDebug()
	}

	//Define new database
	{
		// Define new database
		dbConfig, err = config.NewDBConfig()
		if err != nil {
			appLogger.Fatalf("Unable to read database configuration: %v", err)
		}

		dbConnector, err = storage.InitGormDB(dbConfig)
		if err != nil {
			appLogger.Fatalf("Unable to establish database: %v", err)
		}

		log.Println("INFO: applying db migration")
		if err := storage.MigrateGorm(dbConnector); err != nil {
			log.Fatalf("FATAL: db migration failed %v", err)
		}
	}

	// Define Compliance Repository
	complianceRepository = storage.InitComplianceRepository(dbConnector, appLogger)

	//Define service for interact with different compliance providers
	externalComplianceService = service.InitExternalComplianceService(appLogger)

	//Define service-client for interact with AO-service
	accountOperatorServiceClient = service.InitAccountOperatorServiceClient(appLogger)

	// Define Compliance ComplianceService
	complianceService = service.InitComplianceService(appLogger, complianceRepository,
		externalComplianceService, accountOperatorServiceClient)

	//todo delete all REST layer. Use only gRPC
	// Define REST API and its Router
	//Define main Controllers
	complianceController = rest.InitComplianceController(appLogger, complianceService)

	return nil
}

func startGRPCServer(errCh chan<- error) (shutdown, error) {
	grpcSrv := grpc.InitGrpcServer(appConfig.GrpcHost, appConfig.GrpcPort, appLogger, complianceService)
	go func() {
		if err := grpcSrv.Run(errCh); err != nil {
			appLogger.Fatalf("gRPC Server failed. Error: '%v'\n", err)
		}
	}()

	return func() error {
		grpcSrv.GracefulStop()
		return nil
	}, nil
}

func startHTTPServer(errCh chan<- error) (shutdown, error) {
	complianceRouter := router.NewRouterHandler(appLogger, *appConfig, complianceController)
	appRouter := complianceRouter.InitRouter()

	//init HTTP httpServer
	httpServer := &fasthttp.Server{
		Handler: appRouter.Handler,
	}

	listeningAddress := net.JoinHostPort(appConfig.Host, appConfig.Port)
	log.Printf("INFO: starting HTTP httpServer listening %v\n", listeningAddress)
	go func() {
		errCh <- httpServer.ListenAndServe(net.JoinHostPort(appConfig.Host, appConfig.Port))
	}()

	return httpServer.Shutdown, nil
}

func waitServers(errCh <-chan error, shutdowns ...shutdown) {
	done := make(chan struct{})
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case err := <-errCh:
			log.Printf("ERROR: got an error: %v\n", err)
		case s := <-signalCh:
			log.Printf("INFO: got a signal: %v\n", s)
		}
		log.Println("INFO: gracefully shutting down...")
		for _, sd := range shutdowns {
			if err := sd(); err != nil {
				log.Printf("ERROR: error during shutdown: %v\n", err)
			}
		}
		close(done)
	}()
	<-done
}
