package logic

import (
	"pb"
	"xcore/xlog"
	"xcore/xrpc"
)

// rpc game server service
type ESService struct {
}

// echo
func (*ESService) Echo(pipe xrpc.Pipe, arg *pb.S2S_Echo) {
	if arg == nil {
		xlog.Errorf("rpc OtherPlaceLogin err, invalid args")
		return
	}

	xlog.Debugf("rpc Echo called! %s", arg.Str)
}
