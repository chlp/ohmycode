package store

import (
	"context"
	"fmt"
	"log"
	"ohmycode_api/pkg/util"
	"reflect"
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

func (db *Db) Select(collection string, filter map[string]interface{}, resultType interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	coll := db.db.Collection(collection)
	cursor, err := coll.Find(ctx, bson.M(filter))
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %v", err)
	}
	defer cursor.Close(ctx)

	sliceType := reflect.SliceOf(reflect.TypeOf(resultType).Elem())
	sliceValue := reflect.MakeSlice(sliceType, 0, 0)
	for cursor.Next(ctx) {
		elem := reflect.New(reflect.TypeOf(resultType).Elem()).Interface()
		if err = cursor.Decode(elem); err != nil {
			return nil, err
		}
		sliceValue = reflect.Append(sliceValue, reflect.ValueOf(elem).Elem())
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed with cursor on select: %v", err)
	}

	return sliceValue.Interface(), nil
}

func (db *Db) Upsert(collection string, document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	coll := db.db.Collection(collection)
	docBytes, err := bson.Marshal(document)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %v", err)
	}

	var docMap bson.M
	err = bson.Unmarshal(docBytes, &docMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal document: %v", err)
	}

	id, ok := docMap["_id"]
	if !ok {
		return fmt.Errorf("document must have an _id field")
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": docMap}

	_, err = coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("failed to upsert document: %v", err)
	}

	return nil
}
