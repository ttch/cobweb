package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"code.google.com/p/go.exp/fsnotify"
)

var conf = flag.String("c", "config.json", "gossamer config file path")
var logger chan string = make(chan string)

func main() {
	flag.Parse()
	buf, err := ioutil.ReadFile(*conf)
	CheckErr(err)
	data := map[string]interface{}{}
	err = json.Unmarshal(buf, &data)

	if gossamer, ok := data["watch"]; ok {
		if configs, ok := gossamer.([]interface{}); ok {
			for _, config := range configs {
				buffer, err := json.Marshal(config)
				CheckErr(err)
				var conf Config
				err = json.Unmarshal(buffer, &conf)
				CheckErr(err)
				go watchIt(conf)
			}
		} else {
			panic(fmt.Errorf("except path settings but got %v\n", configs))
		}
	} else {
		panic(fmt.Errorf("Which paths not found in conf[\"watch\"]\n"))
	}
	for {
		message := <-logger
		fmt.Println(message)
	}
}

type Config struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Path    string      `json:"path"`
	Action  []string    `json:"action"`
	Include interface{} `json:"include"`
	Exclude interface{} `json:"exclude"`
}

func watchIt(conf Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Errorf("watcher create failed: %v", err))
	}
	defer func() {
		watcher.Close()
		logger <- fmt.Sprintf("watcher on %s closed.", conf.Path)
	}()

	logger <- fmt.Sprintf("watching %s ...", conf.Path)
	err = filepath.Walk(conf.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(fmt.Errorf("walk to %s and got error %v", path, err))
		}
		if info.IsDir() {
			err = watcher.Watch(path)
			if err != nil {
				panic(fmt.Errorf("watcher can't watch %v , got: %v", path, err))
			}

		}
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("watcher can't watch %v , got: %v", conf.Path, err))
	}

	for {
		select {
		case ev := <-watcher.Event:
			OnNotify(conf, ev)
		case err := <-watcher.Error:
			panic(fmt.Errorf("got error: %v when watch %v", err, conf.Name))
		}
	}

}

func OnNotify(conf Config, event *fsnotify.FileEvent) {
	if conf.Include != nil {
		include := conf.Include.(string)
		ok, err := filepath.Match(include, filepath.Base(event.Name))
		CheckErr(err)
		if !ok {
			return
		}
	}
	if conf.Exclude != nil {
		exclude := conf.Exclude.(string)
		ok, err := filepath.Match(exclude, filepath.Base(event.Name))
		CheckErr(err)
		if ok {
			return
		}
	}
	logger <- fmt.Sprintf("watcher %v : %v", conf.Name, event)
	RunCommand(conf.Action)
}

func RunCommand(command []string) {
	fmt.Println(command)
	cmd := exec.Command(command[0], command[1:]...)
	message, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	if message != nil {
		logger <- (string)(message)
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
