package lib

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var mongoStore *MongoStore
var onceSync sync.Once

func GetMongoStore() *MongoStore {
	onceSync.Do(func() {
		mongoStore = &MongoStore{}
	})
	return mongoStore
}

type MongoStore struct {
	mgoCli *mongo.Client
	DB     *mongo.Database
}

// CreateConn 和数据库连接连接
func (store *MongoStore) CreateConn(address string, dataBaseName string) error {
	clientOptions := options.Client().ApplyURI(address)
	clientOptions.SetMaxConnecting(50)
	clientOptions.SetMaxConnIdleTime(60 * time.Minute)
	mgoCli, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	// 检查连接
	err = mgoCli.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	store.mgoCli = mgoCli
	dataBase := mgoCli.Database(dataBaseName)
	store.DB = dataBase
	return nil
}

// FindOne 查询单个
func (store *MongoStore) FindOne(ctx context.Context, collectionName string, filter bson.D, decodeModel any) error {
	collection := store.DB.Collection(collectionName)
	store.mgoCli.Disconnect(ctx)
	queryResult := collection.FindOne(nil, filter)
	if queryResult.Err() == mongo.ErrNoDocuments {
		return errors.New("NotFound")
	}
	err := queryResult.Decode(decodeModel)
	if err != nil {
		return err
	}
	return nil
}
