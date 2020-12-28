package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Department - struct.
type Department struct {
	ID        string    `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Repository - interface.
type Repository interface {
	Create(ctx context.Context, user *Department) error
}

// MongoRepository - struct.
type MongoRepository struct {
	mongo *mongo.Collection
}

// NewRepository - returns MongoRepository pointer.
func NewRepository(mongo *mongo.Collection) *MongoRepository {
	return &MongoRepository{mongo}
}
