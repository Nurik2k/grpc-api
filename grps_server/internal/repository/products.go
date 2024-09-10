package repository

import (
	"context"
	"fmt"

	product "github.com/BalamutDiana/grps_server/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Products struct {
	db *mongo.Database
}

func NewProducts(db *mongo.Database) *Products {
	return &Products{
		db: db,
	}
}

func (r *Products) Insert(ctx context.Context, item product.Product) error {
	_, err := r.db.Collection("products").InsertOne(ctx, item)

	return err
}

func (r *Products) GetByName(ctx context.Context, name string) (product.Product, error) {
	var prod product.Product
	filter := bson.D{{Key: "name", Value: name}}

	err := r.db.Collection("products").FindOne(ctx, filter).Decode(&prod)
	if err != nil {
		return prod, err
	}
	return prod, nil
}

func (r *Products) UpdateByName(ctx context.Context, prod product.Product) error {
	filter := bson.D{{Key: "name", Value: prod.Name}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "price", Value: prod.Price},
		}},
		{Key: "$inc", Value: bson.D{
			{Key: "changes_count", Value: 1},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "timestamp", Value: prod.Timestamp},
		}},
	}

	_, err := r.db.Collection("products").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Products) List(ctx context.Context, paging product.PagingParams, sorting product.SortingParams) ([]product.Product, error) {
	opts := options.Find()
	sortOpts := bson.D{
		{Key: fmt.Sprintf("%v", sorting.Field), Value: sorting.Asc},
	}

	opts.SetSort(sortOpts)
	opts.SetSkip(int64(paging.Offset))
	opts.SetLimit(int64(paging.Limit))

	cur, err := r.db.Collection("products").Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var productsList []product.Product
	for cur.Next(ctx) {
		var elem product.Product
		if err := cur.Decode(&elem); err != nil {
			return nil, err
		}
		productsList = append(productsList, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return productsList, nil
}
