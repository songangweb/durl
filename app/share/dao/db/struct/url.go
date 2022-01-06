package dbstruct

import (
	"durl/app/share/tool"
	"github.com/xormplus/builder"
	"github.com/xormplus/xorm"
)

type UrlStruct struct {
	Id             int    `xorm:" int pk notnull autoincr"`
	ShortNum       int    `xorm:" int notnull index"`
	FullUrl        string `xorm:" varchar(2048) default('') notnull"`
	ExpirationTime int    `xorm:" int notnull default(0)"`
	IsDel          int    `xorm:" int default(0) notnull"`
	IsFrozen       int    `xorm:" int default(0) notnull"`
	CreateTime     int    `xorm:" created int default(0) notnull"`
	UpdateTime     int    `xorm:" updated int default(0) notnull"`
}

func (I *UrlStruct) TableName() string {
	return "durl_url"
}

const (
	UrlIsDelYes = 1
	UrlIsDelNo  = 0
)

// 函数名称: GetFullUrlByShortNum
// 功能: 通过 ShortNum 获取 完整连接
// 输入参数:
//		ShortNum
// 输出参数:
//		urlDetail: UrlStruct
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func GetFullUrlByShortNum(engine *xorm.EngineGroup, shortNum int) (*UrlStruct, error) {
	urlDetail := new(UrlStruct)

	has, err := engine.
		Where(" short_num = ? and is_del = ? and (expiration_time > ? or expiration_time = ?)",
			shortNum, UrlIsDelNo, tool.TimeNowUnix(), UrlIsDelNo).Get(urlDetail)
	if nil != err {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return urlDetail, nil
}

// 函数名称: GetCacheUrlAllByLimit
// 功能: 查询出符合条件的limit条url
// 输入参数:
//		limit
// 输出参数:
//		urlList: []UrlStruct
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func GetCacheUrlAllByLimit(engine *xorm.EngineGroup, limit int) (*[]UrlStruct, error) {
	urlList := make([]UrlStruct, 0)

	err := engine.
		Where(" is_del = ? and (expiration_time > ? or expiration_time = ?) ",
			UrlIsDelNo, tool.TimeNowUnix(), UrlIsDelNo).
		Limit(limit, 0).
		Find(&urlList)
	return &urlList, err
}

// 函数名称: InsertUrlOne
// 功能: 插入一条数据
// 输入参数:
//		urlStructReq: UrlStruct
// 输出参数:
//		affected: id
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func InsertUrlOne(engine *xorm.EngineGroup, urlStructReq UrlStruct) (int, error) {
	url := new(UrlStruct)
	url.FullUrl = urlStructReq.FullUrl
	url.ShortNum = urlStructReq.ShortNum
	url.ExpirationTime = urlStructReq.ExpirationTime
	affected, err := engine.Insert(url)
	return int(affected), err
}

// 函数名称: DelUrlByShortNum
// 功能: 通过shortNum删除数据
// 输入参数:
//		shortNum: shortNum
// 输出参数:
// 返回: 操作结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func DelUrlByShortNum(engine *xorm.EngineGroup, shortNum int) (bool, error) {

	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	urlStruct := new(UrlStruct)
	urlStruct.IsDel = UrlIsDelYes
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
//	   shortNum: 短链num
// 输出参数:
// 返回: 操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/24 5:14 下午 #

func DelUrlById(engine *xorm.EngineGroup, id int, shortNum int) (bool, error) {

	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	urlStruct := new(UrlStruct)
	urlStruct.IsDel = UrlIsDelYes
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

// 函数名称: UpdateUrlByShortNum
// 功能: 通过shortNum修改数据
// 输入参数:
//		shortNum: 短链num
//		data
// 输出参数:
// 返回: 操作结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func UpdateUrlByShortNum(engine *xorm.EngineGroup, shortNum int, data *map[string]interface{}) (bool, error) {

	session := engine.NewSession()
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
		if key == "shortUrl" {
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

// 函数名称: UpdateUrlById
// 功能: 通过Id修改数据
// 输入参数:
//		id
//		data
// 输出参数:
// 返回: 操作结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 20:44 下午 #

func UpdateUrlById(engine *xorm.EngineGroup, id int, shortNum int, data *map[string]interface{}) (bool, error) {
	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	_, err = session.Table(new(UrlStruct)).Where("id = ?", id).Update(data)
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

func GetShortUrlList(engine *xorm.EngineGroup, fields map[string]interface{}, page, size int) (*[]UrlStruct, error) {
	urlList := make([]UrlStruct, 0)

	q := engine.Where("is_del = ? ", UrlIsDelNo)
	if fields["shortKey"] != nil {
		// 格式转换
		shortNum := tool.Base62Decode(fields["shortKey"].(string))
		q.And(builder.Eq{"short_num": shortNum})
	}
	if fields["fullUrl"] != nil {
		q.And(builder.Like{"full_url", fields["fullUrl"].(string)})
	}
	if fields["isFrozen"] != nil {
		q.And(builder.Eq{"is_frozen": fields["isFrozen"]})
	}
	if fields["createTimeL"] != nil {
		q.And(builder.Gte{"create_time": fields["createTimeL"]})
	}
	if fields["createTimeR"] != nil {
		q.And(builder.Lte{"create_time": fields["createTimeR"]})
	}

	err := q.Limit(size, (page-1)*size).Desc("create_time").Find(&urlList)
	return &urlList, err
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

func GetShortUrlListTotal(engine *xorm.EngineGroup, fields map[string]interface{}) (int, error) {
	urlCount := new(UrlStruct)

	q := engine.Where("is_del = ? ", UrlIsDelNo)
	if fields["fullUrl"] != nil {
		q.And(builder.Like{"full_url", fields["fullUrl"].(string)})
	}
	if fields["shortKey"] != nil {
		// 格式转换
		shortNum := tool.Base62Decode(fields["shortKey"].(string))
		q.And(builder.Eq{"short_num": shortNum})
	}
	if fields["shortNum"] != nil {
		// 格式转换
		q.And(builder.Eq{"short_num": fields["shortNum"]})
	}
	if fields["isFrozen"] != nil {
		q.And(builder.Eq{"is_frozen": fields["isFrozen"]})
	}
	if fields["createTimeL"] != nil {
		q.And(builder.Gte{"create_time": fields["createTimeL"]})
	}
	if fields["createTimeR"] != nil {
		q.And(builder.Lte{"create_time": fields["createTimeR"]})
	}

	total, err := q.Count(urlCount)
	return int(total), err
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

func GetShortUrlInfo(engine *xorm.EngineGroup, fields map[string]interface{}) (*UrlStruct, error) {
	urlDetail := new(UrlStruct)

	q := engine.Where("is_del = ? ", UrlIsDelNo)
	if fields["id"] != nil {
		q.And(builder.Eq{"id": fields["id"]})
	}
	_, err := q.Get(urlDetail)
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

func GetAllShortUrl(engine *xorm.EngineGroup, fields map[string]interface{}) (*[]UrlStruct, error) {
	urlList := make([]UrlStruct, 0)

	q := engine.Where("is_del = ? ", UrlIsDelNo)
	if fields["id"] != nil {
		q.And(builder.Eq{"id": fields["id"]})
	}
	err := q.Find(&urlList)
	return &urlList, err
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

func BatchUpdateUrlByIds(engine *xorm.EngineGroup, updateWhere map[string]interface{}, insertShortNum []int, updateData *map[string]interface{}) (bool, error) {

	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	q := session.Table(new(UrlStruct)).Where("is_del = ?", UrlIsDelNo)
	if updateWhere["id"] != nil {
		q.And(builder.Eq{"id": updateWhere["id"]})
	}
	_, err = q.Update(updateData)
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
