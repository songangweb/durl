package xormDb

import (
	comm "durl/app/share/comm"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
)

type MysqlConf struct {
	Master     string
	Slave1     string
	SetMaxOpen int
	SetMaxIdle int
}

func InitMysql(m MysqlConf) {

	var err error

	// 主库
	master, err := xorm.NewEngine("mysql", m.Master)
	if err != nil {
		defer fmt.Println(comm.MsgCheckDbMysqlConf)
		panic(comm.MsgDbMysqlConfError + ", err: " + fmt.Errorf("%v", err).Error())
	}

	// 从库,如有多个增加相关配置及代码即可
	slave1, err := xorm.NewEngine("mysql", m.Slave1)
	if err != nil {
		defer fmt.Println(comm.MsgCheckDbMysqlConf)
		panic(comm.MsgDbMysqlConfError + ", err: " + fmt.Errorf("%v", err).Error())
	}

	slaves := []*xorm.Engine{slave1}
	Engine, err = xorm.NewEngineGroup(master, slaves)
	if err != nil {
		defer fmt.Println(comm.MsgCheckDbMysqlConf)
		panic(comm.MsgDbMysqlConfError + ", err: " + fmt.Errorf("%v", err).Error())
	}

	// 判断数据库是否链接成功
	err = Engine.Ping()
	if err != nil {
		defer fmt.Println(comm.MsgCheckDbMysqlConf)
		panic(comm.MsgDbMysqlConfError + ", err: " + fmt.Errorf("%v", err).Error())
	}

	Engine.SetMaxOpenConns(m.SetMaxOpen)
	Engine.SetMaxIdleConns(m.SetMaxIdle)
}
