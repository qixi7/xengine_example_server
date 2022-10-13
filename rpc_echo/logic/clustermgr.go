package logic

import (
	"fmt"
	"github.com/qixi7/xengine_core/xhlink"
	"github.com/qixi7/xengine_core/xnet"
	"github.com/qixi7/xengine_core/xrpc"
	"xserver/rpc_echo/global"
	"xserver/rpc_echo/loader"
	"xserver/rpc_echo/root"
)

func NewClusterMgrHlink() ([]*xhlink.HLinkUnitConfig, []*xhlink.HLinkUnitConfig, error) {
	path := fmt.Sprintf("%s/cluster.json", root.ServerConfigDir)
	selfNode := loader.GetServerEnv().GetNodeName()
	listen, dial, err := global.GetRPCStatic().NewClusterHLinkFromFile(path, selfNode, nil)
	if err != nil {
		return nil, nil, err
	}
	return listen, dial, err
}

// --------------------------- rpc方法 ---------------------------

// call ID
func ClusterCallID(remoteID int32, srvName string, arg xnet.ProtoMessage, cb xrpc.Callback) bool {
	return global.GetRPCStatic().CallID(remoteID, srvName, arg, cb)
}

// call Name
func ClusterCallName(nodeName string, srvName string, arg xnet.ProtoMessage, cb xrpc.Callback) bool {
	return global.GetRPCStatic().CallName(nodeName, srvName, arg, cb)
}
