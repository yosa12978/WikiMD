package mongo

import (
	"context"

	"github.com/yosa12978/WikiMD/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *mongo.Database
)

func InitMongo(cfg *config.Config) (*mongo.Database, error) {
	options := options.Client().ApplyURI(cfg.Mongo.Conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db = client.Database(cfg.Mongo.Db)
	return db, nil
}

func GetClient() *mongo.Database {
	return db
}
