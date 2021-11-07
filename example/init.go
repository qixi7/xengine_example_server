package main

import (
	"example/loader"
	"example/logic"
	"example/root"
	"github.com/json-iterator/go"
	"xcore/xlog"
	"xcore/xmodule"
)

func doInitData() bool {
	root.DataRoot.Register(root.EDServerEnv, &loader.ServerConfig{})
	return root.DataRoot.Load()
}

func doInitGame() bool {
	// game logic about
	root.GameRoot.Register(root.EGAgentMgr, &logic.AgentMgr{})
	return root.GameRoot.InitAll()
}

func doInitStatic() bool {
	root.StaticRoot.Register(root.ESNetThread, NewNetThread())
	return root.StaticRoot.InitAll()
}

func doInit() bool {
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
	serverEnv := loader.GetServerEnv()
	if serverEnv == nil {
		xlog.Errorf("printServerInfo err.")
		return
	}
	jsonCfg := jsoniter.Config{IndentionStep: 4}.Froze()
	cfgStr, _ := jsonCfg.MarshalToString(serverEnv)
	xlog.InfoF("init with args:\n%s", cfgStr)
}
