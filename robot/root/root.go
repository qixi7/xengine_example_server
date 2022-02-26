package root

import (
	"math/rand"
	"xcore/xmodule"
)

// Data Root
var DataRoot xmodule.SModuleMgr

// Game Root
var GameRoot xmodule.DModuleMgr

// Static Root
var StaticRoot xmodule.DModuleMgr

// 常规随机(非伪随机)
var NormalRand *rand.Rand

// EData
const (
	EDServerEnv int = iota // 服务器配置
	EDMax
)

// EGame
const (
	EGRobotMgr int = iota // 连接管理
	EGTimer
	EGMax
)

// EStatic
const (
	ESJob int = iota // job worker
	ESMax
)
