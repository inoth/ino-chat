package db

import (
	"context"
	"inochat/client/config"
	"inochat/client/db/model"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	if err := mogo.Collection(res.Col()).FindOne(context.TODO(), filter).Decode(res); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func FindAll(filter interface{}, res model.IEntity) ([]bson.M, error) {
	ctx := context.TODO()
	cur, err := mogo.Collection(res.Col()).Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	r := make([]bson.M, 0)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var tmp bson.M
		if err = cur.Decode(&tmp); err != nil {
			log.Fatal(err)
		}
		r = append(r, tmp)
	}
	return r, nil
}

func Create(entity model.IEntity) bool {
	if _, err := mogo.Collection(entity.Col()).InsertOne(context.TODO(), entity); err != nil {
		logrus.Errorf("%v", err)
		return false
	}
	return true
}

func UpdateOne(filter interface{}, entity model.IEntity) bool {
	cnt, err := mogo.Collection(entity.Col()).UpdateOne(context.TODO(), filter, entity)
	if err != nil {
		return false
	}
	return cnt.ModifiedCount > 0
}

func UpdateMany(filter interface{}, entity model.IEntity) bool {
	cnt, err := mogo.Collection(entity.Col()).UpdateMany(context.TODO(), filter, entity)
	if err != nil {
		return false
	}
	return cnt.ModifiedCount > 0
}

func ToStruct(bsonValue interface{}, res model.IEntity) error {
	data, err := bson.Marshal(bsonValue)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(data, res)
	if err != nil {
		return err
	}
	return nil
}

func GetDB() *mongo.Database {
	return mogo
}
