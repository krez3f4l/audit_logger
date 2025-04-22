package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/krez3f4l/audit_logger/internal/config"
	"github.com/krez3f4l/audit_logger/internal/repository"
	"github.com/krez3f4l/audit_logger/internal/server"
	"github.com/krez3f4l/audit_logger/internal/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	opts := options.Client()
	opts.SetAuth(options.Credential{
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
	})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.DB.Database)

	auditRepo := repository.NewAudit(db)
	auditService := service.NewService(auditRepo)
	auditSrv := server.NewAuditServer(auditService)
	srv := server.NewServer(auditSrv)

	fmt.Println("SERVER STARTED", time.Now())

	if err := srv.ListenAndServe(cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
