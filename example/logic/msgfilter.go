package logic

import "xserver/pb"

type IMsgReceiver interface {
	Recv(ag *agent, pack interface{})
}

// 消息注册
func (d *MsgDispatcher) registerAllMsgHandler() {
	// 登陆
	d.registerMsgHandler(pb.MSG_SyncRoleInfo, (*pb.C2S_SyncRoleInfo)(nil), &recvLogin{})
}
