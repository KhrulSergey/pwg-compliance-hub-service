package main

import (
	grpc "gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/grpc/server"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/rest"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/http/router"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/service"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/storage"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/config"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/logger"
)

func main() {
	appLogger, _ := logger.NewRelease()
	defer func() {
		_ = appLogger.Flush()
	}() // flushes buffer, if any

	// Define new appConfig
	appConfig, err := config.InitAppConfig()
	if err != nil {
		appLogger.Fatalf("Unable to read configuration %v", err)
	}
	//update logger level
	if appConfig.LoggerMode == "DEBUG" {
		appLogger, _ = logger.NewDebug()
	}

	// Define new database
	var dbConnector *gorm.DB
	{
		// Define new database
		dbConfig, err := config.NewDBConfig()
		if err != nil {
			appLogger.Fatalf("Unable to read database configuration: %v", err)
		}

		dbConnector, err = storage.InitGormDB(dbConfig)
		if err != nil {
			appLogger.Fatalf("Unable to establish database: %v", err)
		}
		if err := storage.MigrateGorm(dbConnector); err != nil {
			appLogger.Fatalf("DB migration failed %v", err)
		}

		connection, err := dbConnector.DB()
		if err != nil {
			appLogger.Fatalf("something wrong with database: %v", err)
		}
		defer connection.Close()
	}

	// Define Compliance Repository
	complianceRepository := storage.InitComplianceRepository(dbConnector, appLogger)

	//Define service for interact with different compliance providers
	externalComplianceService := service.InitExternalComplianceService(appLogger)

	accountOperatorServiceClient := service.InitAccountOperatorServiceClient(appLogger)

	// Define Compliance ComplianceService
	complianceService := service.InitComplianceService(appLogger, complianceRepository, externalComplianceService, accountOperatorServiceClient)

	//todo delete all REST layer. Use only gRPC
	// Define REST API and its Router
	//Define main Controllers
	complianceController := rest.InitComplianceController(appLogger, complianceService)

	complianceRouter := router.InitRouterHandler(appLogger, *appConfig, complianceController)
	appRouter := complianceRouter.InitRouter()

	//init GRPC server
	appLogger.Info("Starting servers...")
	grpcSrv := grpc.InitGrpcServer(appConfig.GrpcHost, appConfig.GrpcPort, appLogger, complianceService)
	go func() {
		if err := grpcSrv.Run(); err != nil {
			appLogger.Fatalf("gRPC Server failed. Error: '%v'\n", err)
		}
	}()

	//init HTTP server
	s := &fasthttp.Server{
		Handler: appRouter.Handler,
	}

	//Run App
	cn := make(chan os.Signal, 1)
	signal.Notify(cn, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cn
		appLogger.Info("gracefully shutting down...")
		_ = s.Shutdown()
		grpcSrv.GracefulStop()
	}()

	appLogger.Info("Starting server...")
	appLogger.Fatal(s.ListenAndServe(net.JoinHostPort(appConfig.Host, appConfig.Port)))
}
