package store

import (
	"context"
	"fmt"
	"log"
	"ohmycode_api/pkg/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConfig struct {
	ConnectionString string          `json:"connectionString"`
	DBName           string          `json:"dbname"`
	Timeout          util.OhDuration `json:"timeout"`
}

type Db struct {
	client  *mongo.Client
	db      *mongo.Database
	timeout time.Duration
}

func newDb(config DBConfig) *Db {
	connectTimeout := config.Timeout.Duration
	if connectTimeout <= 0 {
		connectTimeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.ConnectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	return &Db{
		client:  client,
		db:      client.Database(config.DBName),
		timeout: connectTimeout,
	}
}

func (db *Db) FindOne(collection string, filter map[string]interface{}, out interface{}) (found bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	coll := db.db.Collection(collection)
	if err := coll.FindOne(ctx, bson.M(filter)).Decode(out); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, fmt.Errorf("failed to find one document: %w", err)
	}
	return true, nil
}

func (db *Db) ReplaceOneUpsert(collection string, filter map[string]interface{}, document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	coll := db.db.Collection(collection)
	_, err := coll.ReplaceOne(ctx, bson.M(filter), document, options.Replace().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("failed to replace(upsert) document: %w", err)
	}

	return nil
}

func (db *Db) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}
