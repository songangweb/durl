package xormDbStruct

import (
	"durl/app/share/dao/db/xormDb"
	"durl/app/share/tool"
	_ "github.com/go-sql-driver/mysql"
)

type UrlStruct struct {
	Id             int    `xorm:" int pk notnull autoincr"`
	ShortNum       int    `xorm:" int notnull index"`
	FullUrl        string `xorm:" varchar(2048) default('') notnull"`
	ExpirationTime int    `xorm:" int notnull default(0)"`
	IsDel          int8   `xorm:" tinyint default(0) notnull"`
	IsFrozen       int8   `xorm:" tinyint default(0) notnull"`
	CreateTime     int    `xorm:" created default(0) notnull"`
	UpdateTime     int    `xorm:" updated default(0) notnull"`
}

func (I *UrlStruct) TableName() string {
	return "durl_url"
}

// GetFullUrlByShortNum 通过 ShortNum 获取 完整连接
func GetFullUrlByShortNum(shortNum int) (*UrlStruct, error) {
	urlDetail := new(UrlStruct)
	has, err := xormDb.Engine.
		Where(" short_num = ? and is_del = ? and (expiration_time > ? or expiration_time = ?)",
			shortNum, 0, tool.TimeNowUnix(), 0).Get(urlDetail)
	if nil != err {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return urlDetail, nil
}

// GetCacheUrlAllByLimit 查询出符合条件的limit条url
func GetCacheUrlAllByLimit(limit int) ([]UrlStruct, error) {
	urlList := make([]UrlStruct, 0)
	err := xormDb.Engine.
		Where(" is_del = ? and (expiration_time > ? or expiration_time = ?) ",
			0, tool.TimeNowUnix(), 0).
		Limit(limit, 0).
		Find(&urlList)
	return urlList, err
}

// InsertUrlOne 插入一条数据
func InsertUrlOne(urlStructReq UrlStruct) (int64, error) {
	url := new(UrlStruct)
	url.FullUrl = urlStructReq.FullUrl
	url.ShortNum = urlStructReq.ShortNum
	url.ExpirationTime = urlStructReq.ExpirationTime
	affected, err := xormDb.Engine.Insert(url)
	return affected, err
}

// DelUrlByShortNum 通过shortNum删除数据
func DelUrlByShortNum(shortNum int) (bool, error) {

	session := xormDb.Engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	urlStruct := new(UrlStruct)
	urlStruct.IsDel = 1
	_, err = session.Where(" short_num = ? ", shortNum).Update(urlStruct)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	var QueueOne QueueStruct
	QueueOne.ShortNum = shortNum
	_, err = session.Insert(QueueOne)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	err = session.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// 函数名称: DelUrlById
// 功能: 通过id删除数据
// 输入参数:
//     id: 数据id
//	   shortNum: 短链Key
// 输出参数:
// 返回: 操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/24 5:14 下午 #

func DelUrlById(id string, shortNum int) (bool, error) {

	session := xormDb.Engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	urlStruct := new(UrlStruct)
	urlStruct.IsDel = 1
	_, err = session.Where(" id = ? ", id).Update(urlStruct)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	// 删除的短链推送到处理列表
	var QueueOne QueueStruct
	QueueOne.ShortNum = shortNum
	_, err = session.Insert(QueueOne)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	err = session.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateUrlByShortNum 通过shortNum修改数据
func UpdateUrlByShortNum(shortNum int, data *map[string]interface{}) (bool, error) {

	session := xormDb.Engine.NewSession()
	defer session.Close()

	err := session.Begin()

	dataVal := make(map[string]interface{})
	for key, val := range *data {
		if key == "expirationTime" {
			dataVal["expiration_time"] = val
			continue
		}
		if key == "isFrozen" {
			dataVal["is_frozen"] = val
			continue
		}
		if key == "url" {
			dataVal["full_url"] = val
			continue
		}
		dataVal[key] = val
	}

	// 修改数据
	_, err = session.Table(new(UrlStruct)).Where("short_num = ?", shortNum).Update(&dataVal)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	var QueueOne QueueStruct
	QueueOne.ShortNum = shortNum
	_, err = session.Insert(QueueOne)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	err = session.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateUrlById 通过Id修改数据
func UpdateUrlById(id string, shortNum int, data map[string]interface{}) (bool, error) {

	session := xormDb.Engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	_, err = session.Table(new(UrlStruct)).Where("id = ?", id).Update(&data)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	var QueueOne QueueStruct
	QueueOne.ShortNum = shortNum
	_, err = session.Insert(QueueOne)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	err = session.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// 函数名称: GetShortUrlList
// 功能: 查询出符合条件url列表信息
// 输入参数:
//	   where: 业务搜索条件
//     page: 页码
//	   size: 每页展示条数
// 输出参数: []UrlStruct
// 返回: 符合条件数据
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/22 11:25 上午 #

func GetShortUrlList(where map[string][]interface{}, page, size int) ([]UrlStruct, error) {
	urlList := make([]UrlStruct, 0)
	whereStr, bindValue := getWhereStr(where)
	err := xormDb.Engine.
		Select("id,short_num,full_url,expiration_time,is_frozen,create_time,update_time").
		Where(whereStr, bindValue...).
		Limit(size, (page-1)*size).
		Find(&urlList)
	return urlList, err
}

// 函数名称: GetShortUrlListTotal
// 功能:  查询出符合条件url列表信息条数
// 输入参数:
//     where: 业务搜索条件
// 输出参数:
//	   结果条数
//     error
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/22 11:27 上午 #

func GetShortUrlListTotal(where map[string][]interface{}) (int64, error) {
	urlCount := new(UrlStruct)
	whereStr, bindValue := getWhereStr(where)
	total, err := xormDb.Engine.
		Where(whereStr, bindValue...).
		Count(urlCount)
	return total, err
}

// 函数名称: GetShortUrlInfo
// 功能: 获取单条ShortUrl详情
// 输入参数:
//     where: 检索条件
// 输出参数:
//	   *UrlStruct: Url结构
//	   error
// 返回: 检索结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/24 5:10 下午 #

func GetShortUrlInfo(where map[string][]interface{}) (*UrlStruct, error) {
	urlDetail := new(UrlStruct)
	whereStr, bindValue := getWhereStr(where)
	_, err := xormDb.Engine.
		Where(whereStr, bindValue...).
		Get(urlDetail)
	return urlDetail, err
}

// 函数名称: GetAllShortUrl
// 功能: 获取检索Url信息无条数限制
// 输入参数:
//     where: 检索条件
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/30 6:13 下午 #

func GetAllShortUrl(where map[string][]interface{}) ([]UrlStruct, error) {
	urlList := make([]UrlStruct, 0)
	sql, args := getWhereStr(where)
	err := xormDb.Engine.
		Where(sql, args...).
		Find(&urlList)
	return urlList, err
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

	session := xormDb.Engine.NewSession()
	defer session.Close()

	err := session.Begin()

	sql, args := getWhereStr(updateWhere)
	// 修改数据
	_, err = session.Table(new(UrlStruct)).Where(sql, args...).Update(updateData)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	queue := make([]QueueStruct, len(insertShortNum))
	for k, v := range insertShortNum {
		queue[k].ShortNum = v
	}
	_, err = session.Insert(&queue)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	err = session.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}
