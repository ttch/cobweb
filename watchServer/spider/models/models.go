package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/astaxie/beego"
)

var listener map[string][]string = nil


func Listener() map[string][]string {
	if listener == nil {
		fn := beego.AppConfig.String("listen")
		data, err := ioutil.ReadFile(fn)
		beego.Error(data)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(data, &listener)
		if err != nil {
			panic(err)
		}
	}
	return listener
}

func RunCommand(command []string) (string, error) {
	fmt.Println("---------")
	fmt.Println(command)
	cmd := exec.Command(command[0], command[1:]...)
	fmt.Println(cmd)
	message, err := cmd.Output()
	if err != nil {
		beego.Error(err)
		return "", err
	}
	if message != nil {
		beego.Info(message)
		return string(message), nil
	} else {
		return "", fmt.Errorf("except a message or error but nil")
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
