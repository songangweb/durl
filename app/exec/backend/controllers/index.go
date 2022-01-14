package controllers

// Index 首页入口
func (c *BackendController) Index() {
	c.TplName = "dist/index.html"
}
