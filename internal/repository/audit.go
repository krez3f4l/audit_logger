package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/krez3f4l/audit_logger/pkg/domain/audit"
)

type Audit struct {
	db *mongo.Database
}

func NewAudit(db *mongo.Database) *Audit {
	return &Audit{
		db: db,
	}
}

func (r *Audit) Insert(ctx context.Context, item audit.LogItem) error {
	_, err := r.db.Collection("logs").InsertOne(ctx, item)

	return err
}
