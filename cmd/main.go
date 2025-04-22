package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/krez3f4l/audit_logger/internal/config"
	"github.com/krez3f4l/audit_logger/internal/repository"
	"github.com/krez3f4l/audit_logger/internal/server"
	"github.com/krez3f4l/audit_logger/internal/service"
)

const (
	timeoutPing = 3 * time.Second
)

var (
	configDir  = os.Getenv("AUDIT_CONFIG_DIR")
	configName = os.Getenv("AUDIT_CONFIG_NAME")
)

func main() {
	cfg, err := config.NewConfig(configDir, configName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), timeoutPing)
	defer cancel()

	opts := options.Client().ApplyURI(cfg.DBConn.URI)
	opts.SetAuth(options.Credential{
		Username: cfg.DBConn.Username,
		Password: cfg.DBConn.Password,
	})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.DBConn.Database)

	auditRepo := repository.NewAudit(db)
	auditService := service.NewService(auditRepo)
	auditSrv := server.NewAuditServer(auditService)
	srv := server.NewServer(auditSrv)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = srv.ListenAndServe(cfg.Server.Port); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	fmt.Println("gRPC audit server started", time.Now())

	<-shutdownChan
	fmt.Println("Shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeoutPing)
	defer shutdownCancel()

	srv.GracefulStop()

	if err = client.Disconnect(shutdownCtx); err != nil {
		log.Printf("Mongo disconnect error: %v", err)
	}

	fmt.Println("Audit server stopped")
}
