package server

import (
	"fmt"
	"net"

	"github.com/krez3f4l/audit_logger/pkg/domain/audit"

	"google.golang.org/grpc"
)

type Server struct {
	grpcSrv     *grpc.Server
	auditServer audit.AuditServiceServer
}

func NewServer(auditServer audit.AuditServiceServer) *Server {
	return &Server{
		grpcSrv:     grpc.NewServer(),
		auditServer: auditServer,
	}
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	audit.RegisterAuditServiceServer(s.grpcSrv, s.auditServer)

	if err := s.grpcSrv.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) GracefulStop() {
	s.grpcSrv.GracefulStop()
}
