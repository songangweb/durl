package main

import (
	"durl/app/exec/backend/appInit"

	"github.com/beego/beego/v2/server/web"
)

func main() {

	// 项目初始化
	appInit.Init()

	// 项目启动
	web.Run()

}
