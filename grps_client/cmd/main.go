package main

import (
	"log"

	products "github.com/BalamutDiana/grps_server/pkg/domain"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":5672", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := products.NewProductsServiceClient(conn)

	// response, err := c.Fetch(context.Background(), &products.FetchRequest{
	// 	Url: "http://164.92.251.245:8080/api/v1/products/",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("STATUS:", response.Status)

	response, err := c.List(context.Background(), &products.ListRequest{
		SortField:    products.ListRequest_PRICE,
		SortAsc:      -1,
		PagingOffset: 1,
		PagingLimit:  1,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(response.Product)
}
