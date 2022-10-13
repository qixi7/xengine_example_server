package logic

import (
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_core/xnet"
	"xserver/pb"
)

type robotHandler struct {
	roleID   uint32 // robot
	connType int    // 连接类型
}

func (h *robotHandler) OnOpen(s xnet.Sessioner, b bool) {
	r := GetRobotMgr().GetRobotByRoleID(h.roleID)
	if r == nil {
		xlog.Errorf("robotHandler OnOpen err, roleID=%d robot not exist", h.roleID)
		return
	}
	r.GetNetM().SetConnected(h.connType, true)
	xlog.Debugf("robotHandler OnOpen. connType=%d, reopen=%t", h.connType, b)
	// syncRoleInfo
	proto := &pb.C2S_SyncRoleInfo{
		RoleID: r.baseM.GetRoleID(),
		Name:   r.baseM.GetName(),
	}
	r.netM.SendGsMsg(pb.MSG_SyncRoleInfo, proto)
}

func (h *robotHandler) OnClose(s xnet.Sessioner) {
	r := GetRobotMgr().GetRobotByRoleID(h.roleID)
	if r == nil {
		//xlog.Errorf("robotHandler OnClose err, roleID=%d robot not exist", h.roleID)
		return
	}
	xlog.Debugf("robotHandler OnClose, connType=%d", h.connType)

	r.GetNetM().SetConnected(h.connType, false)
	GetRobotMgr().DelRobot(h.roleID, "OnClose")
	r.netM.CloseConnect(h.connType)
}

func (h *robotHandler) OnMessage(s xnet.Sessioner, bindata []byte, pk *xnet.PBPacket) {
	if pk == nil {
		return
	}
	handler, ok := msgHandlerMgr[pk.MsgID]
	if !ok {
		xlog.Errorf("robot handler, msgID=%d, not register", pk.MsgID)
		s.Close()
		return
	}
	r := GetRobotMgr().GetRobotByRoleID(h.roleID)
	if r == nil {
		//xlog.Errorf("robotHandler OnMessage err, roleID=%d robot not exist", h.roleID)
		return
	}
	handler.Recv(r, pk.Body)
}
