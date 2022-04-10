package logic

import (
	"xcore/xlog"
	"xcore/xnet"
)

// ----------- MS Server Handler ------------
type ExampleHandler struct {
}

func (h *ExampleHandler) OnOpen(node string, link *xnet.Link) {
	xlog.InfoF("Example node=%s, remoteID=%d OnOpen", node, link.GetRemoteID())
}

func (h *ExampleHandler) OnClose(node string, link *xnet.Link) {
	xlog.InfoF("Example=%s, remoteID=%d OnClose", node, link.GetRemoteID())
}

func (h *ExampleHandler) OnShutdown(node string, remote int32) {
	xlog.InfoF("Example node=%s, remoteID=%d OnShutdown", node, remote)
}

// ----------- other Server Handler ------------
