package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/aclgo/product/config"
	"github.com/aclgo/product/internal/product/delivery/grpc/service"
	"github.com/aclgo/product/internal/product/repository"
	"github.com/aclgo/product/internal/product/usecase"
	"github.com/aclgo/product/migrations"
	"github.com/aclgo/product/pkg/postgres"
	"github.com/aclgo/product/proto"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.NewConfig(".")
	if err := cfg.Load(); err != nil {
		log.Fatalf("cfg.Load: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := postgres.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("postgres.Connect: %v", err)
	}

	if cfg.MigrationRunning {
		migrations.SetAppMigratioFs(db, nil)
		if err := migrations.Run(); err != nil {
			log.Fatal(err)
		}
	}

	repo := repository.NewPostgresRepo(db)

	productUC := usecase.NewProductUseCase(repo)

	svc := service.NewserviceGRPC(productUC)

	listen, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	server := grpc.NewServer()

	proto.RegisterProductServiceServer(server, svc)

	log.Printf("grpc running port %s", cfg.ServerPort)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("server.Server: %v", err)
	}
}
