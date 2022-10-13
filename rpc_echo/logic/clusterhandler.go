package logic

import (
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_core/xnet"
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

// ----------- other Server Handler ------------
