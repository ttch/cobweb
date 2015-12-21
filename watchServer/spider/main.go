package main

import (
    "github.com/astaxie/beego"
    "fmt"
	"github.com/ttch/watchser/models"
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

		fmt.Println([]string{oCmd.Command,oCmd.Name})
		
		message, err := models.RunCommand([]string{oCmd.Command,oCmd.Name})
		
		models.CheckErr(err)
		fmt.Println(message)
		this.Ctx.WriteString(message)
	}

}

func main() {
    beego.Router("/listen/:action", &MainController{})
    beego.Run()
}
