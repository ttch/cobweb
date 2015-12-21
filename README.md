#简介


CobWeb主要是为了解决在docker中，如果你的工程文件，比如.go程序，在mac或者pc端更新编辑后，docker内没有做正确的响应，更新文件.

使用方法十分简单。

确保docker内已经安装`beego`后,在docker内执行

	go get github.com/ttch/cobweb
	
然后在go/src/github.com/ttch/cobweb/watchServer/spider目录下执行
	
	bee run

这时候服务启动后。

在mac端，同时安装

	go get github.com/howeyc/fsnotify
	go get github.com/ttch/cobweb

然后在go/src/github.com/ttch/cobweb/watchClient目录执行

	go run TRCC.go

config.json是配置文件，执行前你可以通过配置选项来实现监视配置，本软件支持多目录配置。

------

注意：键盘党和急躁党，用的时候，稍微耐心点，会很爽。