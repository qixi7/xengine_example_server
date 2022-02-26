package main

import (
	"github.com/json-iterator/go"
	"math/rand"
	"robot/loader"
	"robot/logic"
	"robot/root"
	"time"
	"xcore/xcontainer/job"
	"xcore/xcontainer/timer"
	"xcore/xlog"
	"xcore/xmodule"
)

func doInitData() bool {
	root.DataRoot.Register(root.EDServerEnv, &loader.RobotConfig{})
	return root.DataRoot.Load()
}

func doInitGame() bool {
	// game logic about
	root.GameRoot.Register(root.EGTimer, timer.New())
	root.GameRoot.Register(root.EGRobotMgr, &logic.RobotMgr{})

	return root.GameRoot.InitAll()
}

func doInitStatic() bool {
	root.StaticRoot.Register(root.ESJob, job.NewController(1024, 4))
	// todo. to be add service

	return root.StaticRoot.InitAll()
}

func doInit() bool {
	root.NormalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	root.DataRoot = xmodule.NewSModuleMgr(root.EDMax)
	root.GameRoot = xmodule.NewDModuleMgr(root.EGMax)
	root.StaticRoot = xmodule.NewDModuleMgr(root.ESMax)
	if !doInitData() {
		return false
	}
	if !doInitGame() {
		return false
	}
	if !doInitStatic() {
		return false
	}

	return true
}

func doDestroy() {
	root.GameRoot.DestroyAll()
	root.StaticRoot.DestroyAll()
}

func printServerInfo() {
	serverEnv := loader.GetRobotCfg()
	if serverEnv == nil {
		xlog.Errorf("printServerInfo err.")
		return
	}
	jsonCfg := jsoniter.Config{IndentionStep: 4}.Froze()
	cfgStr, _ := jsonCfg.MarshalToString(serverEnv)
	xlog.InfoF("init with args:\n%s", cfgStr)
}
