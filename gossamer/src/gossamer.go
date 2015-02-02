package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"code.google.com/p/go.exp/fsnotify"
)

var conf = flag.String("c", "config.json", "gossamer config file path")
var logger chan string = make(chan string)
var count = 0

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
				watchIt(conf)
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
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Path        string      `json:"path"`
	Action      []string    `json:"action"`
	Include     interface{} `json:"include"`
	Exclude     interface{} `json:"exclude"`
	ExcludeDirs interface{} `json:"exclude-dirs"`
}

func watchIt(conf Config) {
	fmt.Printf("watching %v ... \n", conf.Path)
	go watchDir(conf, conf.Path)
}

func watchDir(conf Config, path string) {
	watcher, err := fsnotify.NewWatcher()
	CheckErr(err)
	defer watcher.Close()
	err = watcher.Watch(path)
	if err != nil {
		panic(fmt.Errorf("watcher can't watch %v , got: %v \n count: %v", path, err, count))
	}
	count += 1
	subs, err := ioutil.ReadDir(path)
	CheckErr(err)
	for _, sub := range subs {
		if sub.IsDir() {
			matchExclude := false
			if conf.ExcludeDirs != nil {
				extds := conf.ExcludeDirs.([]interface{})
				for _, extd := range extds {
					extpath := extd.(string)
					matchExclude, err = filepath.Match(extpath, sub.Name())
					CheckErr(err)
					if matchExclude {
						break
					}
				}
			}
			subpath := filepath.Join(path, sub.Name())
			if !matchExclude {
				go watchDir(conf, subpath)
			}
		}
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
