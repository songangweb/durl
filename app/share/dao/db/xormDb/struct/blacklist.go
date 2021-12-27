package xormDbStruct

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db/xormDb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/builder"
)

type BlacklistStruct struct {
	Id         int    `xorm:" int pk notnull autoincr"`
	Ip         string `xorm:" varchar(2048) default('') notnull"`
	IsDel      int8   `xorm:" tinyint default(0) notnull"`
	CreateTime int    `xorm:" created default(0) notnull"`
	UpdateTime int    `xorm:" updated default(0) notnull"`
}

func (I *BlacklistStruct) TableName() string {
	return "durl_blacklist"
}

// InsertBlacklistOne 插入一条数据
func InsertBlacklistOne(urlStructReq BlacklistStruct) (int64, error) {
	url := new(BlacklistStruct)
	url.Ip = urlStructReq.Ip
	affected, err := xormDb.Engine.Insert(url)
	return affected, err
}

// 函数名称: DelBlacklistById
// 功能: 通过id删除数据
// 输入参数:
//     id: 数据id
// 输出参数:
// 返回: 操作结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func DelBlacklistById(id string) (bool, error) {

	session := xormDb.Engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	blacklistStruct := new(BlacklistStruct)
	blacklistStruct.IsDel = 1
	_, err = session.Where(" id = ? ", id).Update(blacklistStruct)
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

// UpdateBlacklistById 通过Id修改数据
func UpdateBlacklistById(id string, data map[string]interface{}) (bool, error) {

	session := xormDb.Engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	_, err = session.Table(new(BlacklistStruct)).Where("id = ?", id).Update(&data)
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

// 函数名称: GetBlacklistList
// 功能: 查询出符合条件黑名单列表信息
// 输入参数:
//	   where: 业务搜索条件
//     page: 页码
//	   size: 每页展示条数
// 输出参数: []BlacklistStruct
// 返回: 符合条件数据
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func GetBlacklistList(fields map[string]interface{}, page, size int) ([]BlacklistStruct, error) {
	BlacklistList := make([]BlacklistStruct, 0)

	q := xormDb.Engine.Where("is_del = ? ", comm.FalseDel)

	if fields["ip"] != nil {
		q.And(builder.Like{"ip", fields["ip"].(string)})
	}

	if fields["createTimeL"] != nil {
		q.And(builder.Lte{"create_time": fields["createTimeL"]})
	}

	if fields["createTimeR"] != nil {
		q.And(builder.Gte{"create_time": fields["createTimeR"]})
	}

	err := q.Limit(size, (page-1)*size).Desc("create_time").Find(&BlacklistList)
	return BlacklistList, err
}

// 函数名称: GetBlacklistListTotal
// 功能:  查询出符合条件黑名单列表信息条数
// 输入参数:
//     where: 业务搜索条件
// 输出参数:
//	   结果条数
//     error
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func GetBlacklistListTotal(fields map[string]interface{}) (int64, error) {
	BlacklistCount := new(BlacklistStruct)

	q := xormDb.Engine.Where("is_del = ? ", comm.FalseDel)

	if fields["fullUrl"] != nil {
		q.And(builder.Like{"full_url", fields["fullUrl"].(string)})
	}

	if fields["createTimeL"] != nil {
		q.And(builder.Lte{"create_time": fields["createTimeL"]})
	}

	if fields["createTimeR"] != nil {
		q.And(builder.Gte{"create_time": fields["createTimeR"]})
	}

	total, err := q.Count(BlacklistCount)
	return total, err
}

// 函数名称: GetBlacklistInfo
// 功能: 获取单条黑名单详情
// 输入参数:
//     where: 检索条件
// 输出参数:
//	   *BlacklistStruct: Blacklist结构
//	   error
// 返回: 检索结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #

func GetBlacklistInfo(fields map[string]interface{}) (*BlacklistStruct, error) {
	blacklistDetail := new(BlacklistStruct)

	q := xormDb.Engine.Where("is_del = ? ", comm.FalseDel)
	if fields["id"] != nil {
		q.And(builder.Eq{"id": fields["id"]})
	}
	_, err := q.Get(blacklistDetail)
	return blacklistDetail, err
}

// 函数名称: GetBlacklistAll
// 功能: 查询出符合条件的黑名单列表
// 输入参数:
//     where: 检索条件
// 输出参数:
//	   *BlacklistStruct: Blacklist结构
//	   error
// 返回: 检索结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/12 5:44 下午 #

func GetBlacklistAll() ([]BlacklistStruct, error) {
	blacklistList := make([]BlacklistStruct, 0)
	err := xormDb.Engine.
		Where(" is_del = ?",
			0).
		Find(&blacklistList)
	return blacklistList, err
}
