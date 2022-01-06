package controllers

import (
	"fmt"
)

// Index 首页入口
func (c *BackendController) Index() {
	fmt.Println("1")
	c.TplName = "dist/index.html"
}
