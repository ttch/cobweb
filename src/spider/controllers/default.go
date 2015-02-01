package controllers

import (
	"spider/models"

	"github.com/astaxie/beego"
)

type ListenController struct {
	beego.Controller
}

func (c *ListenController) Get() {
	action := c.Ctx.Input.Param(":action")
	if action != "" {
		listener := models.Listener()
		if act, ok := listener[action]; ok {
			message, err := models.RunCommand(act)
			models.CheckErr(err)
			c.Data["json"] = map[string]interface{}{
				"message": message,
				"error":   err}
			c.ServeJson()
		}
	}
	c.Abort("401")
}
