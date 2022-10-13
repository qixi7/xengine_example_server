package logic

import (
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_core/xnet"
	"xserver/pb"
)

// ----------- MS Server Handler ------------
type EchoHandler struct {
}

func (h *EchoHandler) OnOpen(node string, link *xnet.Link) {
	xlog.InfoF("Echo node=%s, remoteID=%d OnOpen", node, link.GetRemoteID())
	proto := &pb.S2S_Echo{
		Str: "Hello~",
	}
	ClusterCallName("rpc_echo_505002", "ESService.Echo", proto, nil)
}

func (h *EchoHandler) OnClose(node string, link *xnet.Link) {
	xlog.InfoF("Echo node=%s, remoteID=%d OnClose", node, link.GetRemoteID())
}

// ----------- other Server Handler ------------
