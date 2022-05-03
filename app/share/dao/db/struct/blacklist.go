package dbstruct

import (
	"github.com/xormplus/builder"
	"github.com/xormplus/xorm"
)

type BlacklistStruct struct {
	Id         int    `xorm:" int pk notnull autoincr"`
	Ip         string `xorm:" varchar(2048) default('') notnull"`
	IsDel      int    `xorm:" int default(0) notnull"`
	CreateTime int    `xorm:" created int default(0) notnull"`
	UpdateTime int    `xorm:" updated int default(0) notnull"`
}

func (I *BlacklistStruct) TableName() string {
	return "durl_blacklist"
}

const (
	BlacklistIsDelYes = 1
	BlacklistIsDelNo  = 0
)

// InsertBlacklistOne
// 函数名称: InsertBlacklistOne
// 功能: 插入一条数据
// 输入参数:
//		urlStructReq: BlacklistStruct
// 输出参数:
// 返回: bool: 操作结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/11/07 20:44 下午 #
func InsertBlacklistOne(engine *xorm.EngineGroup, urlStructReq *BlacklistStruct) (int, error) {
	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()
	url := new(BlacklistStruct)
	url.Ip = urlStructReq.Ip
	affected, err := engine.Insert(url)
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	// 添加消息
	var QueueOne QueueStruct
	QueueOne.QueueType = QueueTypeBlacklistAdd
	QueueOne.Data = urlStructReq.Ip
	_, err = session.Insert(QueueOne)
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}

	return int(affected), nil
}

// DelBlacklistById
// 函数名称: DelBlacklistById
// 功能: 通过id删除数据
// 输入参数:
//     id: 数据id
// 输出参数:
// 返回: bool: 操作结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/11/07 20:44 下午 #
func DelBlacklistById(engine *xorm.EngineGroup, id int) (bool, error) {

	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 修改数据
	blacklistStruct := new(BlacklistStruct)
	blacklistStruct.IsDel = BlacklistIsDelYes
	_, err = session.Where(" id = ? ", id).Update(blacklistStruct)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}

	// 查询ip
	blacklistDetail := new(BlacklistStruct)
	q := engine.Where(" id = ? ", id)
	_, err = q.Get(blacklistDetail)

	// 添加消息
	var QueueOne QueueStruct
	QueueOne.QueueType = QueueTypeBlacklistDel
	QueueOne.Data = blacklistDetail.Ip
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

// UpdateBlacklistById
// 函数名称: UpdateBlacklistById
// 功能: 通过Id修改数据
// 输入参数:
//     id: 数据id
//     data: 修改内容
// 输出参数:
//	   *BlacklistStruct: Blacklist结构
//	   error
// 返回: 检索结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/11/07 20:44 下午 #
func UpdateBlacklistById(engine *xorm.EngineGroup, id int, data *map[string]interface{}) (bool, error) {

	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()

	// 查询ip
	blacklistDetail := new(BlacklistStruct)
	q := engine.Where(" id = ? ", id)
	_, err = q.Get(blacklistDetail)

	datac := *data
	ip := datac["ip"]
	if blacklistDetail.Ip != ip {
		// 添加消息
		var QueueOne QueueStruct
		QueueOne.QueueType = QueueTypeBlacklistAdd
		QueueOne.Data = ip.(string)
		_, err = session.Insert(QueueOne)
		if err != nil {
			_ = session.Rollback()
			return false, err
		}
		QueueOne.QueueType = QueueTypeBlacklistDel
		QueueOne.Data = blacklistDetail.Ip
		_, err = session.Insert(QueueOne)
		if err != nil {
			_ = session.Rollback()
			return false, err
		}
	}

	// 修改数据
	_, err = session.Table(new(BlacklistStruct)).Where("id = ?", id).Update(data)
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

// GetBlacklistList
// 函数名称: GetBlacklistList
// 功能: 查询出符合条件黑名单列表信息
// 输入参数:
//	   where: 业务搜索条件
//		- ip: ip
//		- createTimeL: 创建时间开始
//		- createTimeR: 创建时间结束
//     page: 页码
//	   size: 每页展示条数
// 输出参数: []BlacklistStruct
// 返回: 符合条件数据
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/11/07 20:44 下午 #
func GetBlacklistList(engine *xorm.EngineGroup, fields map[string]interface{}, page, size int) (*[]BlacklistStruct, error) {
	BlacklistList := make([]BlacklistStruct, 0)

	q := engine.Where("is_del = ? ", BlacklistIsDelNo)

	if fields["ip"] != nil {
		q.And(builder.Like{"ip", fields["ip"].(string)})
	}

	if fields["createTimeL"] != nil {
		q.And(builder.Gte{"create_time": fields["createTimeL"]})
	}

	if fields["createTimeR"] != nil {
		q.And(builder.Lte{"create_time": fields["createTimeR"]})
	}

	err := q.Limit(size, (page-1)*size).Desc("create_time").Find(&BlacklistList)
	return &BlacklistList, err
}

// GetBlacklistListTotal
// 函数名称: GetBlacklistListTotal
// 功能:  查询出符合条件黑名单列表信息条数
// 输入参数:
//     where: 业务搜索条件
//		- ip: ip
//		- createTimeL: 创建时间开始
//		- createTimeR: 创建时间结束
// 输出参数:
//	   total: 结果条数
//     error
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/11/07 20:44 下午 #
func GetBlacklistListTotal(engine *xorm.EngineGroup, fields map[string]interface{}) (int, error) {
	BlacklistCount := new(BlacklistStruct)

	q := engine.Where("is_del = ? ", BlacklistIsDelNo)

	if fields["ip"] != nil {
		q.And(builder.Like{"ip", fields["ip"].(string)})
	}

	if fields["createTimeL"] != nil {
		q.And(builder.Gte{"create_time": fields["createTimeL"]})
	}

	if fields["createTimeR"] != nil {
		q.And(builder.Lte{"create_time": fields["createTimeR"]})
	}

	total, err := q.Count(BlacklistCount)
	return int(total), err
}

// GetBlacklistInfo
// 函数名称: GetBlacklistInfo
// 功能: 获取单条黑名单详情
// 输入参数:
//     id: 数据id
// 输出参数:
//	   *BlacklistStruct: Blacklist结构
//	   error
// 返回: 检索结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/11/07 20:44 下午 #
func GetBlacklistInfo(engine *xorm.EngineGroup, fields map[string]interface{}) (*BlacklistStruct, error) {
	blacklistDetail := new(BlacklistStruct)

	q := engine.Where("is_del = ? ", BlacklistIsDelNo)

	if fields["id"] != nil {
		q.And(builder.Eq{"id": fields["id"]})
	}

	_, err := q.Get(blacklistDetail)
	return blacklistDetail, err
}

// GetBlacklistAll
// 函数名称: GetBlacklistAll
// 功能: 查询出符合条件的黑名单列表
// 输入参数:
// 输出参数:
//	   []BlacklistStruct: Blacklist结构
//	   error
// 返回: 检索结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/11/07 20:44 下午 #
func GetBlacklistAll(engine *xorm.EngineGroup) (*[]BlacklistStruct, error) {
	blacklistList := make([]BlacklistStruct, 0)
	err := engine.
		Where(" is_del = ?", BlacklistIsDelNo).Find(&blacklistList)
	return &blacklistList, err
}
