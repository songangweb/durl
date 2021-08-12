package mongoDbStruct

import (
	"context"
	"durl/app/share/dao/db/mongoDb"
	"durl/app/share/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueueStruct struct {
	Id         primitive.ObjectID `bson:"_id"`
	ShortNum   int                `bson:"short_num"`
	IsDel      int                `bson:"is_del"`
	CreateTime int64              `bson:"create_time"`
	UpdateTime int64              `bson:"update_time"`
}

func (I *QueueStruct) TableName() string {
	return "durl_queue"
}

// InsertQueueOne 插入一条数据
func InsertQueueOne(req QueueStruct) (interface{}, error) {
	Q := new(QueueStruct)
	collection := mongoDb.Engine.Collection(Q.TableName())
	Q.Id = primitive.NewObjectID()
	Q.ShortNum = req.ShortNum
	timeUnix := tool.TimeNowUnix()
	Q.CreateTime = timeUnix
	Q.UpdateTime = timeUnix
	insertResult, err := collection.InsertOne(context.Background(), Q)
	if err != nil {
		return "", nil
	}
	return insertResult.InsertedID, nil
}

// ReturnQueueLastId 获取最新一条数据的id
func ReturnQueueLastId() (interface{}, error) {
	var QueueDetail QueueStruct
	collection, err := mongoDb.Engine.Collection(QueueDetail.TableName()).Clone()
	if collection == nil {
		return nil, err
	}
	filter := bson.D{{"is_del", 0}}
	findOptions := options.FindOne().SetSort(bson.D{{"_id", -1}})
	err = collection.FindOne(context.Background(), filter, findOptions).Decode(&QueueDetail)
	if err != nil {
		return "", nil
	}
	return QueueDetail.Id, err
}

// GetQueueListById 获取需要处理的数据
func GetQueueListById(id interface{}) ([]*QueueStruct, error) {

	var all []*QueueStruct
	var err error
	var Q QueueStruct
	collection, err := mongoDb.Engine.Collection(Q.TableName()).Clone()
	if collection == nil {
		return nil, err
	}
	filter := bson.D{{"_id", bson.D{{"$gte", id}}}, {"is_del", 0}}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return all, err
	}
	if err = cur.Err(); err != nil {
		return all, err
	}
	err = cur.All(context.Background(), &all)
	if err != nil {
		return all, err
	}
	err = cur.Close(context.Background())
	if err != nil {
		return nil, err
	}

	return all, err
}
