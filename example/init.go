package main

import (
	"github.com/json-iterator/go"
	"github.com/qixi7/xengine_core/xhlink"
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_core/xmodule"
	"github.com/qixi7/xengine_core/xrpc"
	"xserver/example/loader"
	"xserver/example/logic"
	"xserver/example/root"
)

func addClusterHandler(rpcdata *xrpc.RPCStatic) bool {
	rpcdata.AddClusterHandler("rpc_echo", &logic.EchoHandler{})
	return true
}

func doInitData() bool {
	root.DataRoot.Register(root.EDServerEnv, &loader.ServerConfig{})
	return root.DataRoot.Load()
}

func doInitGame() bool {
	// game logic about
	root.GameRoot.Register(root.EGAgentMgr, &logic.AgentMgr{})
	// rpc
	rpcData := root.GameRoot.Register(root.EGRpcData, xrpc.NewRPCDynamicData())
	rpcStatic := xrpc.NewRPCStatic(rpcData)
	root.StaticRoot.Register(root.ESCluster, rpcStatic)
	if !addClusterHandler(rpcStatic) {
		return false
	}
	return root.GameRoot.InitAll()
}

func doInitStatic() bool {
	root.StaticRoot.Register(root.ESNetThread, NewNetThread())
	// rpc
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
