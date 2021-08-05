package log

import (
	"durl/comm"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
)

type Conf struct {
	Conf string
}

func (c Conf) InitLog() {

	// 初始化
	logs.GetLogger()
	err := logs.SetLogger(logs.AdapterMultiFile, c.Conf)

	if err != nil {
		defer fmt.Println(comm.MsgCheckLogConf)
		panic(comm.MsgLogConfIsError + ", err: " + fmt.Errorf("%v", err).Error())
	}

	// 输入行号
	logs.EnableFuncCallDepth(true)

	// 异步输入日志
	logs.Async()

	//异步输出允许设置缓冲 chan 的大小
	logs.Async(1e3)

	//fmt.Println("log")
	//logs.Debug("my book is bought in the year of ", 2016)
	//logs.Info("this %s cat is %v years old", "yellow", 3)
	//logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
	//logs.Error(1024, "is a very", "good game")
	//logs.Critical("oh,crash")
}
