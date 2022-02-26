package main

import (
	"flag"
	"robot/global"
	"robot/loader"
	"robot/root"
	"xcore/xlog"
	"xcore/xutil"
)

// 进程启动参数
var serverParams = struct {
}{}

func init() {
}

func main() {
	flag.Parse()
	// you can set log level. if you like.
	xutil.GApp.SetFPS(global.DefaultFPS).SetVersion(0, 1, 0)
	if !xutil.GApp.Init(appInit) {
		return
	}
	printServerInfo()
	xlog.LogLevel = loader.GetRobotCfg().LogLevel
	xutil.GApp.Run(mainLoop)
	xutil.GApp.Destroy(destroy)
}

func appInit() bool {
	return doInit()
}

func mainLoop(dt int64) {
	// for test
	root.StaticRoot.RunAll(1)
	root.GameRoot.RunAll(1)
}

func destroy() {
	doDestroy()
}
