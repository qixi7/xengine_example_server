package logic

import (
	"github.com/qixi7/xengine_core/xcontainer/timer"
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_core/xmodule"
	"github.com/qixi7/xengine_core/xutil"
	"runtime"
	"time"
	"xserver/robot/global"
	"xserver/robot/loader"
	"xserver/robot/root"
)

/*
	robotmgr.go: 管理所有机器人
*/

type RobotMgr struct {
	baseRobotRoleID uint32               // 基础robotID, 用于生产roleID, 代替http服务器
	robotMap        map[uint32]*oneRobot // key->roleID, value->robot
	genRobotNum     int                  // 生产机器人总数
}

// --------------------------- 内置函数 ---------------------------

// 生成roleID
func (s *RobotMgr) genRoleID() uint32 {
	s.baseRobotRoleID++
	return s.baseRobotRoleID
}

// --------------------------- 外置函数 ---------------------------

// get
func GetRobotMgr() *RobotMgr {
	return root.GameRoot.Getter(root.EGRobotMgr).Get().(*RobotMgr)
}

// 根据roleID获取robot
func (s *RobotMgr) GetRobotByRoleID(roleID uint32) *oneRobot {
	return s.robotMap[roleID]
}

// 获取当前有多少机器人
func (s *RobotMgr) GetNowRobotNum() int {
	return len(s.robotMap)
}

// 获取当前gs、bs连接数情况.
// 返回: gsConnNum, bsConnNum
func (s *RobotMgr) GetConnNum() (int, int) {
	gsConnNum, bsConnNum := 0, 0
	for rid := range s.robotMap {
		netM := s.robotMap[rid].GetNetM()
		if netM.GetConnected(global.EConnectTypeGS) {
			gsConnNum++
		}
		if netM.GetConnected(global.EConnectTypeBS) {
			bsConnNum++
		}
	}
	return gsConnNum, bsConnNum
}

func (s *RobotMgr) DelRobot(roleID uint32, reason string) {
	r, ok := s.robotMap[roleID]
	if !ok {
		return
	}
	r.netM.CloseConnect(global.EConnectTypeGS)
	r.netM.CloseConnect(global.EConnectTypeBS)
	s.robotMap[roleID] = nil
	s.genRobotNum--
	delete(s.robotMap, roleID)
	xlog.InfoF("DelRobot roleID=%d, reason=%s", roleID, reason)
}

// --------------------------- impl interface ---------------------------

// Init impl DModule
func (s *RobotMgr) Init(selfGetter xmodule.DModuleGetter) bool {
	s.baseRobotRoleID = 1000000
	s.robotMap = make(map[uint32]*oneRobot)
	// 辅助测试: 如果内存过大, 直接关闭
	global.GetTimerMgr().Add(time.Second*10, &checkMemTimer{})
	return true
}

// Run impl DModule
func (s *RobotMgr) Run(delta int64) {
	cfg := loader.GetRobotCfg()
	nowNum := s.genRobotNum
	if nowNum < cfg.ClientNum {
		// 计算该帧实际连多少个client
		createNum := cfg.ClientNum - nowNum
		if createNum > cfg.ConnNumPerFrame {
			createNum = cfg.ConnNumPerFrame
		}
		// go connect.
		for i := 0; i < createNum; i++ {
			// new robot
			roleID := s.genRoleID()
			newRobot := NewRobot(roleID)
			s.robotMap[roleID] = newRobot
			// go connect gs
			connType := global.EConnectTypeGS
			newRobot.netM.ConnectGS(&robotHandler{roleID: roleID, connType: connType})
			s.genRobotNum++
			xlog.Debugf("gen new Robot roleID=%d", roleID)
		}
	}
	// robot update
	for rid := range s.robotMap {
		s.robotMap[rid].netM.Update()
	}
}

// Destroy impl DModule
func (s *RobotMgr) Destroy() {
}

// 检查内存情况timer
type checkMemTimer struct {
}

func (t *checkMemTimer) Invoke(ctl *timer.Controller, item *timer.Item) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	maxMemInUse := uint64(6 * 1024 * 1024 * 1024) // 4G
	if m.HeapInuse >= maxMemInUse {
		xlog.InfoF("HeapInuse=%d kb, over max...\n", m.HeapInuse/1024)
		xutil.GApp.Exit()
		return
	}
	ctl.AddItem(item)
}
