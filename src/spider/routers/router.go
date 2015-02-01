package routers

import (
	"spider/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/listen/:action", &controllers.ListenController{})
}
