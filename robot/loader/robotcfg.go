package loader

import (
	"fmt"
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_pub/cfgloader"
	"xserver/robot/root"
)

// 机器人配置
type RobotConfig struct {
	ProjName        string // 项目名
	NodeType        string // 节点服务器类型
	WorldID         int    // 世界ID
	ServerIP        string // 服务器IP地址
	GsTcpPort       int    // gs登陆需要的tcp端口
	ClientNum       int    // 客户端数量
	ConnNumPerFrame int    // 每帧连多少个客户端
	LogLevel        int    // 日志层级
	// 额外信息
	Program string // nodeName+RemoteID组成的标识
}

func (cfg *RobotConfig) Path() string {
	return "./config/robotconfig.json"
}

// --------------------------- impl reload interface ---------------------------

func (cfg *RobotConfig) Load() bool {
	if !cfgloader.LoadJsonFile(cfg) {
		return false
	}
	return cfg.checkValid()
}

func (cfg *RobotConfig) Reload() {
}

func (cfg *RobotConfig) Destroy() {
	// do nothing
}

// --------------------------- health check ---------------------------

// 健康检查,顺便可以做些额外的事, 比如你可以缓存一些数据到struct, 方便后续操作.
func (cfg *RobotConfig) checkValid() bool {
	cfg.Program = fmt.Sprintf("%s_%s_%d", cfg.NodeType, cfg.ProjName, cfg.WorldID)
	return true
}

func (cfg *RobotConfig) PrintRobotConfig() {
	xlog.Debugf("PrintRobotConfig begin.")
	xlog.Debugf("ProjName=%s, WorldID=%d, Program=%s",
		cfg.ProjName, cfg.WorldID, cfg.Program)
	xlog.Debugf("ServerIP=%s, ClientNum=%d, ConnNumPerFrame=%d",
		cfg.ServerIP, cfg.ClientNum, cfg.ConnNumPerFrame)
	xlog.Debugf("PrintRobotConfig end.")
}

// ------------- All Check is done. begin your show -----------------

// Get 配置
func GetRobotCfg() *RobotConfig {
	return root.DataRoot[root.EDServerEnv].(*RobotConfig)
}
