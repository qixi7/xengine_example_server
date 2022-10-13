package logic

import (
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_core/xmodule"
	"github.com/qixi7/xengine_core/xnet"
	"xserver/example/root"
	"xserver/pb"
)

// 客户端管理类
type AgentMgr struct {
	key2agent    map[AgentKey]*agent // key->agentKey, value->agent
	roleID2agent map[uint32]*agent   // key->roleID, value->agent
	dispatchMgr  *MsgDispatcher      // 消息分发管理类
}

// --------------------------- 内置函数 ---------------------------

// --------------------------- 外置函数 ---------------------------

// get module
func GetAgentMgr() *AgentMgr {
	return root.GameRoot.Getter(root.EGAgentMgr).Get().(*AgentMgr)
}

// 通过roleID获取agent
func (c *AgentMgr) GetAgentByRoleID(roleID uint32) *agent {
	ag, ok := c.roleID2agent[roleID]
	if !ok {
		return nil
	}
	return ag
}

// 客户端连接成功回调
func (c *AgentMgr) OnAgentConnect(remoteID int32, s xnet.Sessioner, clientIP string) {
	// new agent
	agentKey := AgentKey{ConnID: s.ID()}
	if _, ok := c.key2agent[agentKey]; ok {
		// 理论上这个if不可能进
		return
	}
	// add link.
	c.key2agent[agentKey] = newAgent(clientIP, agentKey, s)
}

// 客户端同步roleID
func (c *AgentMgr) SyncRoleInfo(ag *agent, arg *pb.C2S_SyncRoleInfo) {
	if ag == nil {
		return
	}
	// add link.
	c.roleID2agent[arg.RoleID] = ag
}

// 客户端断开连接回调
func (c *AgentMgr) OnAgentDisConnect(remoteID int32, conn uint32) {
	key := AgentKey{
		ConnID: conn,
	}
	ag, ok := c.key2agent[key]
	if !ok {
		//xlog.Errorf("OnAgentDisConnect err, key=%v, not online", key)
		return
	}
	c.DirectDestroyAgent(ag.agentKey)
}

func (c *AgentMgr) DirectDestroyAgent(key AgentKey) {
	// 延迟掉线后调用, 一定时间后销毁agent
	ag, ok := c.key2agent[key]
	if !ok {
		//xlog.Errorf("DirectDestroyAgent err, key=%v, not online", key)
		return
	}
	xlog.InfoF("AgentMgr: DirectDestroyAgent, key=%v", key)
	delete(c.key2agent, key)
	delete(c.roleID2agent, ag.GetRoleID())
	ag.session.Close()
	ag = nil
}

// 分发消息
func (c *AgentMgr) DispatchMsg(remoteID int32, packet *xnet.PBPacket) {
	key := AgentKey{ConnID: packet.ConnID}
	ag, ok := c.key2agent[key]
	if !ok {
		xlog.Errorf("DispatchMsg err, agentKey=%v, not exist", key)
		return
	}
	c.dispatchMgr.dispatch(ag, packet)
}

// --------------------------- impl interface ---------------------------

// Init impl DModule
func (c *AgentMgr) Init(selfGetter xmodule.DModuleGetter) bool {
	c.roleID2agent = make(map[uint32]*agent)
	c.key2agent = make(map[AgentKey]*agent)
	c.dispatchMgr = NewMsgDispatcher()
	return true
}

// Run impl DModule
func (c *AgentMgr) Run(delta int64) {
}

// Destroy impl DModule
func (c *AgentMgr) Destroy() {
}
