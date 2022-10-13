package main

import (
	"github.com/json-iterator/go"
	"github.com/qixi7/xengine_core/xcontainer/job"
	"github.com/qixi7/xengine_core/xcontainer/timer"
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_core/xmodule"
	"math/rand"
	"time"
	"xserver/robot/loader"
	"xserver/robot/logic"
	"xserver/robot/root"
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
