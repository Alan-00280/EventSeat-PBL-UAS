package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NewMongoDatabase membuat koneksi ke MongoDB dan mengembalikan instance *mongo.Database.
// Client disimpan di dalam MongoConn agar bisa di-Disconnect saat shutdown.
type MongoConn struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDatabase menginisialisasi koneksi MongoDB menggunakan URI dari config.Env.
func NewMongoDatabase(env *config.Env) (*MongoConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(env.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("database: gagal connect mongodb: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("database: mongodb ping gagal: %w", err)
	}

	return &MongoConn{
		Client:   client,
		Database: client.Database(env.MongoDBName),
	}, nil
}