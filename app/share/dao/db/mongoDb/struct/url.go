package mongoDbStruct

import (
	"context"
	"durl/app/share/dao/db/mongoDb"
	"durl/app/share/tool"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UrlStruct struct {
	Id             primitive.ObjectID  `bson:"_id"`
	ShortNum       int                 `bson:"short_num"`
	FullUrl        string              `bson:"full_url"`
	ExpirationTime int                 `bson:"expiration_time"`
	IsDel          int8                `bson:"is_del"`
	IsFrozen       int8                `bson:"is_frozen"`
	CreateTime     primitive.Timestamp `bson:"create_time"`
	UpdateTime     primitive.Timestamp `bson:"update_time"`
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

// 函数名称: DelUrlById
// 功能: 通过Id删除数据
// 输入参数:
//     id: 数据id
//     shortNum: 短链Key
// 输出参数:
// 返回: 操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/24 5:16 下午 #

func DelUrlById(id string, shortNum int) (bool, error) {

	var urlDetail UrlStruct

	collection, err := mongoDb.Engine.Collection(urlDetail.TableName()).Clone()
	if collection == nil {
		return false, err
	}

	// id 需要转换下
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	// 修改url表数据
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"is_del": 1}, "$currentDate": bson.M{"update_time": bson.M{"$type": "timestamp"}}}
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

// 函数名称: UpdateUrlById
// 功能: 通过Id修改数据
// 输入参数:
//     id: 数据id
//     shortNum: 短链key
//     data: 修改的内容
// 输出参数:
// 返回: 操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/25 6:41 下午 #

func UpdateUrlById(id string, shortNum int, data map[string]interface{}) (bool, error) {

	var UrlDetail UrlStruct

	// 修改url表数据
	collection, err := mongoDb.Engine.Collection(UrlDetail.TableName()).Clone()
	if collection == nil {
		return false, nil
	}

	// id 转换
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": data, "$currentDate": bson.M{"update_time": bson.M{"$type": "timestamp"}}}

	err = mongoDb.Engine.Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		_, err = collection.UpdateOne(sessionContext, filter, update)
		if err != nil {
			return err
		}

		// 插入处理列表
		var QueueOne QueueStruct
		QueueOne.ShortNum = shortNum
		_, err = InsertQueueOne(QueueOne)
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		sessionContext.CommitTransaction(sessionContext)
		return nil
	})

	if err != nil {
		return false, nil
	}

	return true, nil
}

// 函数名称:
// 功能: 查询出符合条件的url列表数据
// 输入参数:
//		filter: bson
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/19 3:30 下午 #

func GetShortUrlList(where map[string][]interface{}, page, size int) ([]*UrlStruct, error) {
	var all []*UrlStruct
	var err error

	var U UrlStruct
	filter := getWhereStrMongo(where)
	collection, err := mongoDb.Engine.Collection(U.TableName()).Clone()
	if collection == nil {
		return all, err
	}
	findOptions := options.Find().SetLimit(int64(size)).SetSkip(int64((page - 1) * size))
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

// 函数名称: GetShortUrlListCount
// 功能: 统计出符合条件的url数据量
// 输入参数:
//     filter: bson
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/22 7:33 下午 #

func GetShortUrlListCount(where map[string][]interface{}) (int64, error) {
	var U UrlStruct
	filter := getWhereStrMongo(where)
	collection, err := mongoDb.Engine.Collection(U.TableName()).Clone()
	if collection == nil {
		return 0, err
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 函数名称: GetShortUrlInfo
// 功能: 获取ShortUrl详情
// 输入参数:
//     where: 检索条件
// 输出参数:
//	   *UrlStruct url结构体
//	   error
// 返回: 检索结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/24 5:11 下午 #

func GetShortUrlInfo(where map[string][]interface{}) (*UrlStruct, error) {

	var urlDetail UrlStruct

	collection, err := mongoDb.Engine.Collection(urlDetail.TableName()).Clone()
	if collection == nil {
		return nil, err
	}

	filter := getWhereStrMongo(where)
	err = collection.FindOne(context.Background(), filter).Decode(&urlDetail)
	if err != nil {
		return &urlDetail, err
	}

	return &urlDetail, nil
}

// 函数名称: GetAllShortUrl
// 功能: 获取全部Url信息无条数限制
// 输入参数:
//     where: 检索条件
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/30 6:13 下午 #

func GetAllShortUrl(where map[string][]interface{}) ([]*UrlStruct, error) {
	var all []*UrlStruct
	var err error

	var U UrlStruct
	filter := getWhereStrMongo(where)
	collection, err := mongoDb.Engine.Collection(U.TableName()).Clone()
	if collection == nil {
		return all, err
	}
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
	return all, nil
}

// 函数名称: BatchUpdateUrlByIds
// 功能: 根据UrlId 修改Url信息
// 输入参数:
//		updateWhere: 修改限制条件
//	    insertShortNum: 涉及修改Url的短链Key值
//		updateData: 修改内容
// 输出参数:
// 返回: 修改结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/30 6:19 下午 #

func BatchUpdateUrlByIds(updateWhere map[string][]interface{}, insertShortNum []int, updateData map[string]interface{}) (bool, error) {

	var UrlDetail UrlStruct

	// 修改url表数据
	collection, err := mongoDb.Engine.Collection(UrlDetail.TableName()).Clone()
	if collection == nil {
		return false, nil
	}
	filter := getWhereStrMongo(updateWhere)

	err = mongoDb.Engine.Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		update := bson.M{"$set": updateData, "$currentDate": bson.M{"update_time": bson.M{"$type": "timestamp"}}}

		_, err = collection.UpdateMany(sessionContext, filter, update)
		if err != nil {
			return err
		}

		// 插入处理列表
		var QueueOne QueueStruct
		var queue []interface{}
		for _, v := range insertShortNum {
			QueueOne.Id = primitive.NewObjectID()
			QueueOne.ShortNum = v
			queue = append(queue, QueueOne)
		}
		_, err = InsertQueueMany(queue)
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		sessionContext.CommitTransaction(sessionContext)
		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
