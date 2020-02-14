package controllers

import (
	"bookms/models"
	"github.com/astaxie/beego/logs"
)

type CategoryController struct {
	BaseController
}

func (c *CategoryController) AddCategory() {
	c.TplName="category/add.html"
	if c.Ctx.Input.IsPost() {
		pid, err := c.GetInt("pid", 0)
		if err != nil {
			c.JsonResult(400, "父级分类id不合法")
		}
		category_name := c.GetString("category_name")
		if category_name == "" {
			c.JsonResult(400, "分类名不能为空")
		}
		description := c.GetString("description")
		sort, err := c.GetInt("sort", 0)

		category := models.Category{
			Pid:pid,
			CategoryName:category_name,
			Description:description,
			Sort:sort,
		}
		err = category.AddCategory()
		if err != nil {
			logs.Error("AddCategory: ", err.Error())
			c.JsonResult(500, "添加分类失败")
		}
		c.JsonResult(200, "添加成功")
	}
}
