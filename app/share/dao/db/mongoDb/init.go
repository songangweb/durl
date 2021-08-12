package mongoDb

import (
	"context"
	comm2 "durl/app/share/comm"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Conf struct {
	Type  string
	Mongo MongoConf
}

func InitMongoDb(c Conf) {
	InitMongo(c.Mongo)
}

type MongoConf struct {
	Uri            string
	Db             string
	SetMaxPoolSize int
}

var Engine *mongo.Database

func InitMongo(m MongoConf) {

	var err error
	// 连接数据库
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(m.Uri).SetMaxPoolSize(uint64(m.SetMaxPoolSize)))
	if err != nil {
		defer fmt.Println(comm2.MsgCheckDbMongoConf)
		panic(comm2.MsgDbMongoConfError + ", err: " + fmt.Errorf("%v", err).Error())
	}

	// 判断服务是不是可用
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		defer fmt.Println(comm2.MsgCheckDbMongoConf)
		panic(comm2.MsgDbMongoConfError + ", err: " + fmt.Errorf("%v", err).Error())
	}

	// 获取数据库
	Engine = client.Database(m.Db)

}
