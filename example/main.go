package main

import (
	"example/global"
	"example/loader"
	"example/root"
	"flag"
	"xcore/xlog"
	"xcore/xutil"
)

func main() {
	flag.Parse()
	// you can set log level. if you like.
	xutil.GApp.SetFPS(global.DefaultFPS).SetVersion(0, 1, 0)
	if !xutil.GApp.Init(appInit) {
		return
	}
	printServerInfo()
	xlog.LogLevel = loader.GetServerEnv().GetLogLV()
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
