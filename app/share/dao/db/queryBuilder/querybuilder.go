package query

import (
	"github.com/go-xorm/builder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 函数名称: GetWhereStr
// 功能: 根据业务筛选条件构造orm的where使用
// 输入参数:
//     where: 业务传入的数据筛选条件
// 输出参数:
//	   whereStr: sql字符串
//     bindValue: sql绑定参数
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/22 10:48 上午 #

func GetWhereStr(where map[string][]interface{}) (whereStr string, bindValue []interface{}) {

	if len(where) != 0 {
		// 制做业务搜索条件
		var sql string
		var args []interface{}
		for key, value := range where {
			if value[0].(string) == "like" {
				sql, args, _ = builder.ToSQL(builder.Like{key, value[1].(string)})
			} else if value[0].(string) == "=" || value[0].(string) == "in" {
				sql, args, _ = builder.ToSQL(builder.Eq{key: value[1]})
			}
			whereStr += sql + " and "
			for _, v := range args {
				bindValue = append(bindValue, v)
			}

		}
		whereStr = whereStr[:len(whereStr)-4]
	}
	return whereStr, bindValue
}

// 函数名称: GetWhereStrMongo
// 功能: 根据业务筛选条件构造orm-mongo的fliter使用
// 输入参数:
//     where: 业务传入的数据筛选条件
// 输出参数:
//	   filter: bson
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/22 7:31 下午 #

func GetWhereStrMongo(where map[string][]interface{}) (filter map[string]interface{}) {

	if len(where) != 0 {
		filter = make(map[string]interface{})
		// 制做业务搜索条件
		for key, value := range where {
			// mongo 的id需要处理下
			if key == "id" {
				key = "_id"
				if _, ok := value[1].([]string); ok == true {
					var finalValue []interface{}
					for _, v := range value[1].([]string) {
						ids, _ := primitive.ObjectIDFromHex(v)
						finalValue = append(finalValue, ids)
					}
					value[1] = finalValue
				} else {
					value[1], _ = primitive.ObjectIDFromHex(value[1].(string))
				}
			}

			if value[0].(string) == "like" {
				filter[key] = bson.M{"$regex": value[1], "$options": "i"}
			} else if value[0].(string) == "in" {
				filter[key] = bson.M{"$in": value[1]}
			} else {
				filter[key] = value[1]
			}
		}
	}
	return
}
