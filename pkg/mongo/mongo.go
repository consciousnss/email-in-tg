package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultConnTimeout = 15 * time.Second
	defaultDBName      = "default_db"
)

type Mongo struct {
	connTimeout time.Duration
	dbName      string

	DB *mongo.Database
}

func New(ctx context.Context, url string) (*Mongo, error) {
	m := &Mongo{
		connTimeout: defaultConnTimeout,
		dbName:      defaultDBName,
	}

	clientOptions := options.Client().ApplyURI(url)
	conn, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, m.connTimeout)
	defer cancel()

	err = conn.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	m.DB = conn.Database(m.dbName)

	return m, nil
}
