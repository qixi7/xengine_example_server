package global

import (
	"example/root"
	"xcore/xhlink"
	"xcore/xrpc"
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
