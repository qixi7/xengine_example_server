package logic

import (
	"github.com/qixi7/xengine_core/xnet"
	"xserver/pb"
)

// 客户端唯一key
type AgentKey struct {
	ConnID uint32 // client connectID
}

// 客户端代理, agent == 客户端
type agent struct {
	session  xnet.Sessioner // session
	agentKey AgentKey       // agent key
	roleid   uint32         // roleid
}

// --------------------------- 内置函数 ---------------------------

// new agent
func newAgent(ip string, key AgentKey, s xnet.Sessioner) *agent {
	ag := &agent{
		session:  s,
		agentKey: key,
	}
	return ag
}

// --------------------------- 外置函数 ---------------------------

// 发消息
func (a *agent) SendMsg(msgID pb.Game_Msg, proto xnet.ProtoMessage) {
	newPack := xnet.New_PBPacket()
	newPack.ConnID = a.agentKey.ConnID
	newPack.MsgID = uint16(msgID)
	newPack.Body = proto
	a.session.AsyncSend(newPack)
}

func (a *agent) SetRoleID(roleid uint32) {
	a.roleid = roleid
}

func (a *agent) GetRoleID() uint32 {
	return a.roleid
}
