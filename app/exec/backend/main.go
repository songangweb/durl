package main

import (
	"durl/app/exec/backend/appInit"

	"github.com/beego/beego/v2/server/web"
)

// main
// 函数名称: main
// 功能: 项目初始化启动
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2020/12/07 5:44 下午 #
func main() {

	// 项目初始化
	appInit.Init()

	// 项目启动
	web.Run()
}
