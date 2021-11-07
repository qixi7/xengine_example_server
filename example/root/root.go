package root

import (
	"xcore/xmodule"
)

// Data Root
var DataRoot xmodule.SModuleMgr

// Game Root
var GameRoot xmodule.DModuleMgr

// Static Root
var StaticRoot xmodule.DModuleMgr

// 服务器配置路径
var ServerConfigDir string

// 服务器档数据路径
var ServerDataDir string

// EData
const (
	EDServerEnv  int = iota // 服务器配置
	EDMax
)

// EGame
const (
	EGAgentMgr int = iota // client mgr
	EGMax
)

// EStatic
const (
	ESNetThread       int = iota // 网络线程
	ESMax
)
