package conn

import (
	"Moreover/setting"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

func InitMongo(config *setting.MongoConfig) {
	url := fmt.Sprintf("mongodb://%s:%d/?maxPoolSize=20&w=majority", config.Host, config.Port)
	mongoClient, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	MongoDB = mongoClient.Database("more")
}
