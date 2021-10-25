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
	//url := fmt.Sprintf("mongodb://%s:%s@%s:%d/?maxPoolSize=20&w=majority", config.Username, config.Password, config.Host, config.Port)
	MongoClient, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(url).SetMaxPoolSize(20))
	MongoDB = MongoClient.Database("more")
}
