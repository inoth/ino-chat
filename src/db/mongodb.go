package db

import (
	"context"
	"ino-chat/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
)

const (
	C_DATABASE = ""
)

var (
	DB *MongoDBConnect
)

type MongoDBConnect struct {
	Mogo *mongo.Client
}

func Init() {
	DB = &MongoDBConnect{
		Mogo: setConnect(),
	}
}

func setConnect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 连接池
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Instance().MongoDB.Host).SetMaxPoolSize(20))
	if err != nil {
		log.Error(err)
	}
	return client
}

type mgo struct {
	database   string
	collection string
}

func NewMgo(database, collection string) *mgo {
	return &mgo{
		database,
		collection,
	}
}

func (m *mgo) InsertOne(value interface{}) bool {
	client := DB.Mogo
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return false
	}
	_, err = col.InsertOne(context.TODO(), value)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (m *mgo) InsertMany(values []interface{}) bool {
	client := DB.Mogo
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return false
	}
	_, err = col.InsertMany(context.TODO(), values)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (m *mgo) FindOne(filter interface{}, res interface{}) bool {
	client := DB.Mogo
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return false
	}
	err = col.FindOne(context.TODO(), filter).Decode(res)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (m *mgo) FindMany(skip, limit int64, filter, sort interface{}) ([]bson.M, bool) {
	client := DB.Mogo
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return nil, false
	}
	var findoptions *options.FindOptions
	if limit > 0 {
		findoptions = &options.FindOptions{}
		findoptions.SetSkip(limit * (skip - 1))
		findoptions.SetLimit(limit)
		findoptions.SetSort(sort)
	}
	cur, err := col.Find(context.Background(), filter, findoptions)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	r := make([]bson.M, 0)
	ctx := context.TODO()
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var tmp bson.M
		if err = cur.Decode(&tmp); err != nil {
			log.Fatal(err)
		}
		r = append(r, tmp)
	}
	return r, true
}

func (m *mgo) Count(filter interface{}) int {
	client := DB.Mogo
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return 0
	}
	c, err := col.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0
	}
	return int(c)
}

func (m *mgo) UpdateOne(filter, update interface{}) bool {
	client := DB.Mogo
	col, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		log.Error(err)
		return false
	}
	_, err = col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}
