package service

import (
	"common/lib"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"micro_tmpl/player/playerDefine"
	"micro_tmpl/player/playerModel"
)

type PlayerService struct {
	mongo *lib.MongoStore
}

func (service *PlayerService) SetMongo(mongo *lib.MongoStore) {
	service.mongo = mongo
}

func (service *PlayerService) GetPlayerInfoByName(ctx context.Context, name string) (*playerModel.Player, error) {
	collection := service.mongo.DB.Collection(playerDefine.PlayerTableName)
	filter := bson.D{{"name", name}}
	queryResult := collection.FindOne(ctx, filter)
	if queryResult.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	playerInfo := &playerModel.Player{}
	queryResult.Decode(playerInfo)
	return playerInfo, nil
}

func (service *PlayerService) InsertPlayer(ctx context.Context, playerInfo *playerModel.Player) error {
	collection := service.mongo.DB.Collection(playerDefine.PlayerTableName)
	playerInfo.Id = "10001"
	_, err := collection.InsertOne(ctx, playerInfo)
	if err != nil {
		return err
	}
	return nil
}
