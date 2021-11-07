package logic

import (
	"os"
	"pb"
	"xcore/xlog"
	"xcore/xnet"
)

// 消息分发器
type MsgDispatcher struct {
	msgHandlerMgr map[uint16]IMsgReceiver // key->MsgID, value->handler
}

// new
func NewMsgDispatcher() *MsgDispatcher {
	mgr := &MsgDispatcher{
		msgHandlerMgr: make(map[uint16]IMsgReceiver),
	}
	mgr.registerAllMsgHandler()
	return mgr
}

// --------------------------- 内置函数 ---------------------------

// 注册所有消息handler
func (d *MsgDispatcher) registerMsgHandler(msgID pb.Game_Msg, ptr xnet.ProtoMessage, handler IMsgReceiver) {
	msgID16 := uint16(msgID)
	if _, ok := d.msgHandlerMgr[msgID16]; ok {
		xlog.Errorf("registerMsgHandler err, msgID=%d repeated", msgID16)
		os.Exit(2)
	}
	d.msgHandlerMgr[msgID16] = handler
	xnet.RegisterPBMsgID(uint16(msgID), ptr)
}

// 分发消息
func (d *MsgDispatcher) dispatch(ag *agent, packet *xnet.PBPacket) {
	if ag == nil || packet == nil {
		return
	}
	handler, ok := d.msgHandlerMgr[packet.MsgID]
	if !ok {
		return
	}
	handler.Recv(ag, packet.Body)
}
