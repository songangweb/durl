package xormDbStruct

import (
	"github.com/go-xorm/builder"
)

// 函数名称: getWhereStr
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

func getWhereStr(where map[string][]interface{}) (whereStr string, bindValue []interface{}) {

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
