package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BalamutDiana/grps_server/internal/config"
	"github.com/BalamutDiana/grps_server/internal/repository"
	"github.com/BalamutDiana/grps_server/internal/server"
	"github.com/BalamutDiana/grps_server/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client()
	opts.SetAuth(options.Credential{
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
	})
	opts.ApplyURI(cfg.DB.URI)

	dbClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := dbClient.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	db := dbClient.Database(cfg.DB.Database)

	productsRepo := repository.NewProducts(db)
	productsService := service.NewProduct(productsRepo)

	productSrv := server.NewProductServer(productsService)
	srv := server.New(productSrv)

	fmt.Println("SERVER STARTED", time.Now())

	go func() {
		if err := srv.ListenAndServe(cfg.Server.Port); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Println("SERVER STOPPED", time.Now())
	srv.Stop()
}
