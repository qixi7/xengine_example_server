package logic

import (
	"pb"
	"xcore/xlog"
	"xcore/xnet"
)

// ----------- MS Server Handler ------------
type EchoHandler struct {
}

func (h *EchoHandler) OnOpen(node string, link *xnet.Link) {
	xlog.InfoF("Echo node=%s, remoteID=%d OnOpen", node, link.GetRemoteID())
	proto := &pb.S2S_Echo{
		Str: "Hello~",
	}
	ClusterCallName("rpc_echo_211", "ESService.Echo", proto, nil)
}

func (h *EchoHandler) OnClose(node string, link *xnet.Link) {
	xlog.InfoF("Echo node=%s, remoteID=%d OnClose", node, link.GetRemoteID())
}

func (h *EchoHandler) OnShutdown(node string, remote int32) {
	xlog.InfoF("Echo node=%s, remoteID=%d OnShutdown", node, remote)
}

// ----------- other Server Handler ------------
