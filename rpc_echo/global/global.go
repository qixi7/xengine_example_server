package global

import (
	"github.com/qixi7/xengine_core/xhlink"
	"github.com/qixi7/xengine_core/xrpc"
	"xserver/rpc_echo/root"
)

/*
	global.go: some global function
*/

func GetHLinkMgr() *xhlink.HLinkMgr {
	return root.StaticRoot.Getter(root.ESHLink).Get().(*xhlink.HLinkMgr)
}

func GetRPCStatic() *xrpc.RPCStatic {
	return root.StaticRoot.Getter(root.ESCluster).Get().(*xrpc.RPCStatic)
}
