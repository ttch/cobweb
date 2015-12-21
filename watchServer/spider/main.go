/*
package main

import (
	_ "spider/routers"

	"github.com/astaxie/beego"
)

func init(){
	beego.SetLevel(beego.LevelDebug)
}

func main() {
	beego.Run()
}
*/

package main

import (
    "github.com/astaxie/beego"
    "fmt"
	"github.com/ttch/watchServer/spider/models"
	"encoding/json"
)

/* command struct */
type Command struct {
	Command        string      `json:"Command"`
	Name           string      `json:"Name"`
	Action         string      `json:"Action"`
}


type MainController struct {
    beego.Controller
}

func (this *MainController) Post() {
	action := this.Ctx.Input.Param(":action")

	var oCmd Command

	json.Unmarshal(this.Ctx.Input.RequestBody, &oCmd)

	if action != "" {
		
		//listener := models.Listener()

		//if act, ok := listener[action]; ok {
		fmt.Println([]string{oCmd.Command,oCmd.Name})
		
		message, err := models.RunCommand([]string{oCmd.Command,oCmd.Name})
		
		models.CheckErr(err)
		fmt.Println(message)
		this.Ctx.WriteString(message)
		
		/*
		c.Data["json"] = map[string]interface{}{
			"message": message,
			"error":   err}
		c.ServeJson()
		fmt.Println(c)
		*/
		//}
	}

}

func main() {
    beego.Router("/listen/:action", &MainController{})
    beego.Run()
}
