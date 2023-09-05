package grpc

import (
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/internal/service"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/logger"
	"net"

	grpc "gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/pkg/grpc/compliance_hub_service"
	goGrpc "google.golang.org/grpc"
)

type Server interface {
	Run() error
	GracefulStop()
}

type grpcServer struct {
	grpc.UnimplementedComplianceGrpcServiceServer
	grpcHost               string
	grpcPort               string
	logger                 logger.Logger
	rpcSrv                 *goGrpc.Server
	complianceCheckService service.ComplianceService
}

func InitGrpcServer(grpcHost string, grpcPort string, logger logger.Logger, complianceService service.ComplianceService) Server {
	s := goGrpc.NewServer()
	srv := &grpcServer{
		grpcHost:               grpcHost,
		grpcPort:               grpcPort,
		logger:                 logger,
		rpcSrv:                 s,
		complianceCheckService: complianceService,
	}
	srv.register()
	return srv
}

func (s *grpcServer) register() {
	// TBD in this place we have to call RegisterServiceServer(s, impl)
	// from generated package
	grpc.RegisterComplianceGrpcServiceServer(s.rpcSrv, s)
}

// Run runs grpc grpcServer
func (s *grpcServer) Run() error {

	addr := net.JoinHostPort(s.grpcHost, s.grpcPort)
	s.logger.Infof("Grpc service started on '%s'", addr)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if err := s.rpcSrv.Serve(l); err != nil {
		return err
	}
	return nil
}

// GracefulStop stops the gRPC grpcServer gracefully. It stops the grpcServer from accepting new connections and RPCs and
// blocks until all the pending RPCs are finished.
func (s *grpcServer) GracefulStop() {
	s.rpcSrv.GracefulStop()
}
