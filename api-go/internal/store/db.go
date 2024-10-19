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
	ConnectionString string        `json:"connectionString"`
	DBName           string        `json:"dbname"`
	Timeout          util.Duration `json:"timeout"`
}

type Db struct {
	client  *mongo.Client
	db      *mongo.Database
	timeout time.Duration
}

func NewDb(config DBConfig) *Db {
	clientOptions := options.Client().ApplyURI(config.ConnectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	return &Db{
		client:  client,
		db:      client.Database(config.DBName),
		timeout: config.Timeout.Duration,
	}
}

func (db *Db) Select(collection string, filter map[string]interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	coll := db.db.Collection(collection)
	cursor, err := coll.Find(ctx, bson.M(filter))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	for cursor.Next(ctx) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (db *Db) Exec(collection string, operation string, document interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	coll := db.db.Collection(collection)
	var res *mongo.InsertOneResult
	var err error

	switch operation {
	case "insert":
		res, err = coll.InsertOne(ctx, document)
	default:
		err = fmt.Errorf("unsupported operation: %s", operation)
	}

	return res, err
}
