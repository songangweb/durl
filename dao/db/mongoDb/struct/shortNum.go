package mongoDbStruct

import (
	"context"
	"durl/dao/db/mongoDb"
	"durl/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ShortNumStruct struct {
	Id         primitive.ObjectID `bson:"_id"`
	MaxNum     int                `bson:"max_num"`
	Step       int                `bson:"step"`
	Version    int                `bson:"version"`
	UpdateTime int64              `bson:"update_time"`
}

func (I *ShortNumStruct) TableName() string {
	return "durl_short_num"
}

// ReturnShortNumPeriod 获取号码段
func ReturnShortNumPeriod() (int, int, error) {
	var shortNumDetail ShortNumStruct

	// 获取数据
	collection, err := mongoDb.Engine.Collection(shortNumDetail.TableName()).Clone()
	if collection == nil {
		return 0, 0, err
	}

	filter := bson.D{}
	err = collection.FindOne(context.Background(), filter).Decode(&shortNumDetail)
	if err != nil {
		return 0, 0, err
	}
	// 修改数据
	shortNumDetail.MaxNum += shortNumDetail.Step
	shortNumDetail.UpdateTime = tool.TimeNowUnix()
	update := bson.M{"$set": shortNumDetail}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return 0, 0, err
	}

	return shortNumDetail.Step, shortNumDetail.MaxNum, nil
}

// 插入第一条默认数据
func InsertFirst() error {
	shortNumDetail := ShortNumStruct{
		Id:         primitive.ObjectID{},
		MaxNum:     100,
		Step:       1,
		Version:    1,
		UpdateTime: time.Now().Unix(),
	}
	_, err := mongoDb.Engine.Collection(shortNumDetail.TableName()).InsertOne(context.Background(), shortNumDetail)
	return err
}
