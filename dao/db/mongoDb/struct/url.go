package mongoDbStruct

import (
	"context"
	"durl/dao/db/mongoDb"
	"durl/tool"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UrlStruct struct {
	Id             primitive.ObjectID `bson:"_id"`
	ShortNum       int                `bson:"short_num"`
	FullUrl        string             `bson:"full_url"`
	ExpirationTime int64              `bson:"expiration_time"`
	IsDel          int                `bson:"is_del"`
	CreateTime     int64              `bson:"create_time"`
	UpdateTime     int64              `bson:"update_time"`
}

func (I *UrlStruct) TableName() string {
	return "durl_url"
}

// GetFullUrlByShortNum 通过 ShortNum 获取 完整连接
func GetFullUrlByShortNum(shortNum int) (*UrlStruct, error) {

	var urlDetail UrlStruct
	collection, err := mongoDb.Engine.Collection(urlDetail.TableName()).Clone()
	if collection == nil {
		return nil, err
	}
	filter := bson.D{
		{"is_del", 0},
		{"short_num", shortNum},
		{"$or",
			[]bson.M{
				{"expiration_time": 0},
				{"expiration_time": bson.D{{"$gt", tool.TimeNowUnix()}}},
			},
		},
	}
	err = collection.FindOne(context.Background(), filter).Decode(&urlDetail)
	if err != nil {
		return &urlDetail, err
	}
	return &urlDetail, nil
}

// GetCacheUrlAllByLimit 查询出符合条件的limit条url
func GetCacheUrlAllByLimit(limit int) ([]*UrlStruct, error) {
	var all []*UrlStruct
	var err error

	var U UrlStruct
	collection, err := mongoDb.Engine.Collection(U.TableName()).Clone()
	if collection == nil {
		return all, err
	}
	filter := bson.D{{"is_del", 0},
		{"$or",
			[]bson.M{
				{"expiration_time": 0},
				{"expiration_time": bson.D{{"$gt", tool.TimeNowUnix()}}},
			},
		},
	}
	findOptions := options.Find().SetLimit(int64(limit))
	cur, err := collection.Find(context.Background(), filter, findOptions)
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
	return all, nil
}

// InsertUrlOne 插入一条数据
func InsertUrlOne(urlStructReq UrlStruct) (interface{}, error) {
	u := new(UrlStruct)
	collection := mongoDb.Engine.Collection(u.TableName())
	urlStructReq.Id = primitive.NewObjectID()
	timeUnix := tool.TimeNowUnix()
	urlStructReq.CreateTime = timeUnix
	urlStructReq.UpdateTime = timeUnix
	insertResult, err := collection.InsertOne(context.Background(), urlStructReq)
	if err != nil {
		return "", err
	}
	return insertResult.InsertedID, nil
}

// DelUrlByShortNum 通过shortNum删除数据
func DelUrlByShortNum(shortNum int) (bool, error) {

	var UrlDetail UrlStruct
	// 修改url表数据
	collection, err := mongoDb.Engine.Collection(UrlDetail.TableName()).Clone()
	if collection == nil {
		return false, err
	}
	filter := bson.M{"short_num": shortNum}
	update := bson.M{"$set": bson.M{"is_del": 1}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false, err
	}
	// 插入处理列表
	var QueueOne QueueStruct
	QueueOne.ShortNum = shortNum
	_, err = InsertQueueOne(QueueOne)
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateUrlByShortNum 通过shortNum修改数据
func UpdateUrlByShortNum(shortNum int, data *map[string]interface{}) (bool, error) {

	var UrlDetail UrlStruct
	// 修改url表数据
	collection, err := mongoDb.Engine.Collection(UrlDetail.TableName()).Clone()
	if collection == nil {
		return false, nil
	}

	filter := bson.M{"short_num": shortNum}

	updateData := bson.M{}
	for key, val := range *data {
		updateData[key] = val
	}

	update := bson.M{"$set": updateData}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false, nil
	}

	// 插入处理列表
	var QueueOne QueueStruct
	QueueOne.ShortNum = shortNum
	_, err = InsertQueueOne(QueueOne)
	if err != nil {
		return false, nil
	}

	return true, nil
}
