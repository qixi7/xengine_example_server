package logic

import (
	"fmt"
	"robot/logic/basemodule"
	"robot/logic/netmodule"
)

/*
	onerobot.go: 一个机器人实现
*/

// 一个机器人信息
type oneRobot struct {
	baseM basemodule.BaseInfoModule // 基础信息
	netM  netmodule.NetModule       // 网络模块
}

// --------------------------- 内置函数 ---------------------------

// --------------------------- 外置函数 ---------------------------

func (r *oneRobot) GetBaseM() *basemodule.BaseInfoModule {
	return &r.baseM
}

func (r *oneRobot) GetNetM() *netmodule.NetModule {
	return &r.netM
}

// new
func NewRobot(roleID uint32) *oneRobot {
	robot := &oneRobot{}
	robot.baseM.Init(roleID, fmt.Sprintf("机器人_%d", roleID))
	robot.netM.Init()
	return robot
}
