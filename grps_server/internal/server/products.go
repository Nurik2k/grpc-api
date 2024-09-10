package server

import (
	"context"

	products "github.com/BalamutDiana/grps_server/gen/products"
)

type ProductService interface {
	List(ctx context.Context, req *products.ListRequest) (*products.ListResponse, error)
	Fetch(ctx context.Context, req *products.FetchRequest) (*products.FetchResponse, error)
}

type ProductServer struct {
	service ProductService
}

func NewProductServer(service ProductService) *ProductServer {
	return &ProductServer{
		service: service,
	}
}

func (s *ProductServer) List(ctx context.Context, req *products.ListRequest) (*products.ListResponse, error) {
	return s.service.List(ctx, req)
}

func (s *ProductServer) Fetch(ctx context.Context, req *products.FetchRequest) (*products.FetchResponse, error) {
	return s.service.Fetch(ctx, req)
}
