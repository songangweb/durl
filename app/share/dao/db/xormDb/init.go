package xormDb

import (
	comm "durl/app/share/comm"
	"fmt"
	"github.com/xormplus/xorm"
)

var Engine *xorm.EngineGroup

type Conf struct {
	Type  string
	Mysql MysqlConf
}

func InitXormDb(c Conf) {
	switch c.Type {
	case "mysql":
		InitMysql(c.Mysql)
	default:
		defer fmt.Println(comm.MsgCheckDbType)
		panic(comm.MsgDbTypeError + ", type: " + c.Type)
	}

}
