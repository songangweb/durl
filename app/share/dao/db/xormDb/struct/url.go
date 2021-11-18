package xormDbStruct

import (
	"durl/app/share/dao/db/xormDb"
	"durl/app/share/tool"
	"fmt"
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

// GetShortUrlList 查询出符合条件url列表信息
func GetShortUrlList(where string, page, size int,bindValue ...interface{}) ([]UrlStruct, error) {
	urlList := make([]UrlStruct, 0)
	err := xormDb.Engine.
		Select("id,short_num,full_url,expiration_time,is_frozen,create_time").
		Where(where,bindValue...).
		Limit(size, (page-1)*size).
		Find(&urlList)
	return urlList, err
}

// GetShortUrlListTotal 查询出符合条件url列表信息条数
func GetShortUrlListTotal(where string,bindValue ...interface{}) (int64, error) {
	urlCount := new(UrlStruct)
	total, err := xormDb.Engine.
		Where(where,bindValue...).
		Count(urlCount)
	fmt.Printf("%+v", total)
	return total, err
}
