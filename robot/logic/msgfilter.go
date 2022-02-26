package logic

import (
	"os"
	"pb"
	"xcore/xlog"
	"xcore/xnet"
)

type IMsgReceiver interface {
	Recv(r *oneRobot, pack interface{})
}

var msgHandlerMgr = map[uint16]IMsgReceiver{} // key->MsgID, value->handler

func init() {
	registerAllMsgHandler()
}

// 注册消息handler
func registerMsgHandler(msgID uint16, ptr xnet.ProtoMessage, handler IMsgReceiver) {
	msgID16 := msgID
	if _, ok := msgHandlerMgr[msgID16]; ok {
		xlog.Errorf("registerMsgHandler err, msgID=%d repeated", msgID16)
		os.Exit(2)
	}
	msgHandlerMgr[msgID16] = handler
	xnet.RegisterPBMsgID(msgID, ptr)
}

// 消息注册
func registerAllMsgHandler() {
	// gs
	registerMsgHandler(uint16(pb.MSG_SyncRoleInfo), (*pb.S2C_SyncRoleInfo)(nil), &recvSyncRoleInfo{})
}
