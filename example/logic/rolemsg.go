package logic

import (
	"github.com/qixi7/xengine_core/xlog"
	"xserver/pb"
)

// ------------------------- 测试登陆相关消息 -------------------------

// 登录
type recvLogin struct {
}

func (*recvLogin) Recv(ag *agent, pack interface{}) {
	res, ok := pack.(*pb.C2S_SyncRoleInfo)
	if !ok {
		return
	}
	xlog.Debugf("key=%v, recvLogin, roleID=%d, name=%s", ag.agentKey, res.RoleID, res.Name)
	// 不允许roleID==0
	if res.RoleID == 0 {
		xlog.Errorf("key=%v, recvLogin err, roleID=%d", ag.agentKey, res.RoleID)
		return
	}
	// do something
	GetAgentMgr().SyncRoleInfo(ag, res)
	proto := &pb.C2S_SyncRoleInfo{
		RoleID: res.RoleID,
		Name:   res.Name,
	}
	ag.SendMsg(pb.MSG_SyncRoleInfo, proto)
}
