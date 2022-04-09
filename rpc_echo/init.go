package main

import (
	"github.com/json-iterator/go"
	"rpc_echo/loader"
	"rpc_echo/logic"
	"rpc_echo/root"
	"xcore/xhlink"
	"xcore/xlog"
	"xcore/xmodule"
	"xcore/xrpc"
)

func addClusterHandler(rpcdata *xrpc.RPCStatic) bool {
	rpcdata.AddClusterHandler("example", &logic.ExampleHandler{})
	if err := rpcdata.Register(&logic.ESService{}); err != nil {
		xlog.Errorf("Register rpc MSService err=%v", err)
		return false
	}
	return true
}

func doInitData() bool {
	root.DataRoot.Register(root.EDServerEnv, &loader.ServerConfig{})
	return root.DataRoot.Load()
}

func doInitGame() bool {
	// game logic about

	// rpc about
	rpcData := root.GameRoot.Register(root.EGRpcData, xrpc.NewRPCDynamicData())
	rpcStatic := xrpc.NewRPCStatic(rpcData)
	root.StaticRoot.Register(root.ESCluster, rpcStatic)
	if !addClusterHandler(rpcStatic) {
		return false
	}
	return root.GameRoot.InitAll()
}

func doInitStatic() bool {
	listen, dial, err := logic.NewClusterMgrHlink()
	if err != nil {
		xlog.Errorf("initNetwork err=%v", err)
		return false
	}
	serverEnv := loader.GetServerEnv()
	if serverEnv == nil {
		xlog.Errorf("serverEnv serverEnv not find")
		return false
	}
	linkMgr := xhlink.NewHLinkMgr(xhlink.NewHLinkConfig(
		serverEnv.GetNodeName(),
		serverEnv.RemoteID,
		listen,
		dial))
	root.StaticRoot.Register(root.ESHLink, linkMgr)
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
