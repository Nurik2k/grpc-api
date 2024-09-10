package service

import (
	"context"
	"strconv"
	"time"

	"github.com/BalamutDiana/grps_server/gen/products"
	"github.com/BalamutDiana/grps_server/pkg/csv"
	"github.com/BalamutDiana/grps_server/pkg/domain"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Repository interface {
	Insert(ctx context.Context, item domain.Product) error
	GetByName(ctx context.Context, name string) (domain.Product, error)
	UpdateByName(ctx context.Context, prod domain.Product) error
	List(ctx context.Context, paging domain.PagingParams, sorting domain.SortingParams) ([]domain.Product, error)
}

type Product struct {
	repo Repository
}

func NewProduct(repo Repository) *Product {
	return &Product{
		repo: repo,
	}
}

func (s *Product) List(ctx context.Context, req *products.ListRequest) (*products.ListResponse, error) {
	paging := domain.PagingParams{
		Offset: int(req.GetPagingOffset()),
		Limit:  int(req.GetPagingLimit()),
	}
	sorting := domain.SortingParams{
		Field: req.GetSortField(),
		Asc:   req.GetSortAsc(),
	}

	items, err := s.repo.List(ctx, paging, sorting)
	if err != nil {
		return nil, err
	}
	var sortedProducts []*products.ProductItem

	for _, x := range items {
		var sortedProduct products.ProductItem
		sortedProduct.Name = x.Name
		sortedProduct.Price = int32(x.Price)
		sortedProduct.Count = int32(x.ChangesCount)
		sortedProduct.Timestamp = timestamppb.New(x.Timestamp)
		sortedProducts = append(sortedProducts, &sortedProduct)
	}

	return &products.ListResponse{
		Product: sortedProducts,
	}, nil
}

func (s *Product) Fetch(ctx context.Context, req *products.FetchRequest) (*products.FetchResponse, error) {
	url := req.Url
	data, err := csv.ReadCSVFromUrl(url)
	if err != nil {
		return nil, err
	}

	for idx, row := range data {
		if idx == 0 {
			continue
		}

		name := row[0]
		price, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, err
		}

		var prod domain.Product

		item, err := s.repo.GetByName(ctx, name)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				prod.Name = name
				prod.Price = price
				prod.ChangesCount = 1
				prod.Timestamp = time.Now()

				if err = s.repo.Insert(ctx, prod); err != nil {
					return nil, err
				}
				continue
			} else {
				return nil, err
			}
		}

		if prod.Price != item.Price {
			prod.Price = price
			prod.Timestamp = time.Now()

			if err = s.repo.UpdateByName(ctx, prod); err != nil {
				return nil, err
			}
		}
	}
	return &products.FetchResponse{
		Status: "OK",
	}, nil
}
