package logic

import (
	"pb"
	"xcore/xlog"
)

// ------------------------- 登陆相关消息 -------------------------

// 同步玩家信息
type recvSyncRoleInfo struct {
}

func (*recvSyncRoleInfo) Recv(r *oneRobot, pack interface{}) {
	res := pack.(*pb.S2C_SyncRoleInfo)
	xlog.Debugf("%s, recvSyncRoleInfo, RoleID=%d, Name=%s",
		r.baseM.BaseString(), res.RoleID, res.Name)
	// do something
}
