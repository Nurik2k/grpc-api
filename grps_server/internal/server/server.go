package server

import (
	"fmt"
	"net"

	"github.com/BalamutDiana/grps_server/gen/products"

	"google.golang.org/grpc"
)

type Server struct {
	grpcSrv       *grpc.Server
	productServer products.ProductsServiceServer
}

func New(productServer products.ProductsServiceServer) *Server {
	return &Server{
		grpcSrv:       grpc.NewServer(),
		productServer: productServer,
	}
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	products.RegisterProductsServiceServer(s.grpcSrv, s.productServer)

	if err := s.grpcSrv.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() func() {
	return s.grpcSrv.GracefulStop
}
