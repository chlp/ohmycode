package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"api/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var dbInstance *Db

func GetDb() *Db {
	if dbInstance != nil {
		return dbInstance
	}

	appConf := conf.LoadApiConf().DB
	clientOptions := options.Client().ApplyURI("mongodb://" + appConf.ServerName)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	dbInstance = &Db{
		Client: client,
		DB:     client.Database(appConf.DBName),
	}
	return dbInstance
}

func (db *Db) Select(collection string, filter interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	coll := db.DB.Collection(collection)
	cursor, err := coll.Find(ctx, filter)
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	coll := db.DB.Collection(collection)
	var res *mongo.InsertOneResult
	var err error

	switch operation {
	case "insert":
		res, err = coll.InsertOne(ctx, document)
	// Add more operations as needed
	default:
		err = fmt.Errorf("unsupported operation: %s", operation)
	}

	return res, err
}
