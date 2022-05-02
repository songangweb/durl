package controllers

// Index
// 函数名称: Index
// 功能: 首页入口
// 输入参数:
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 6:41 下午 #
func (c *BackendController) Index() {
	c.TplName = "dist/index.html"
}
