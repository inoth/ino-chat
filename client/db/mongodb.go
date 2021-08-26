package db

import (
	"context"
	"inochat/client/config"
	"inochat/client/db/model"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	mogo *mongo.Database
)

// func Init(dbname string) error {
// 	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Instance().MongoDB))
// 	if err != nil {
// 		return err
// 	}
// 	if err = mongoClient.Ping(context.TODO(), nil); err != nil {
// 		return err
// 	}

// 	Instance().Db = mongoClient.Database(dbname)
// 	return nil
// }

func init() {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Instance().MongoDB))
	if err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
	if err = mongoClient.Ping(context.TODO(), nil); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}

	db := mongoClient.Database("wechat_util")
	mogo = db
}

func GetDb() *mongo.Database {
	return mogo
}

func FindOne(filter interface{}, res model.IEntity) error {
	if err := mogo.Collection(res.ColName()).FindOne(context.TODO(), filter).Decode(res); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func Find(filter interface{}, res model.IEntity) ([]model.IEntity, error) {
	ctx := context.TODO()
	cur, err := mogo.Collection(res.ColName()).Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	r := make([]model.IEntity, 0)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		if err = cur.Decode(&res); err != nil {
			log.Fatal(err)
		}
		r = append(r, res)
	}
	return r, nil
}

func Create(entity model.IEntity) error {
	if _, err := mogo.Collection(entity.ColName()).InsertOne(context.TODO(), entity); err != nil {
		return err
	}
	return nil
}

func GetDB() *mongo.Database {
	return mogo
}
