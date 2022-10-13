package loader

import (
	"fmt"
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_pub/cfgloader"
	"github.com/qixi7/xengine_pub/netutil"
	"strconv"
	"xserver/rpc_echo/root"
)

// 服务器配置
type ServerConfig struct {
	ProjName        string // 项目名
	NodeType        string // 节点服务器类型
	ServerID        int    // 服务器编号
	ZoneID          int    // 同类型服务器区分编号
	LogLevel        int    // 日志层级
	ServerDataDir   string // 服务器档目录
	ServerConfigDir string // 服务器配置目录
	// 额外信息
	RemoteID int32  // 唯一路由ID=ServerID+ZoneID. 比如ServerID=1, ZoneID=2,那么RemoteID=12
	NodeName string // 节点名: NodeType+RemoteID
	SelfIP   string
}

func (cfg *ServerConfig) Path() string {
	return "serverconfig.json"
}

// --------------------------- impl reload interface ---------------------------

func (cfg *ServerConfig) Load() bool {
	if !cfgloader.LoadJsonFile(cfg) {
		return false
	}
	return cfg.checkValid()
}

func (cfg *ServerConfig) Reload() {
}

func (cfg *ServerConfig) Destroy() {
	// do nothing
}

func (cfg *ServerConfig) ReloadName() string {
	return "serverconfig"
}

func (cfg *ServerConfig) ReloadCreate() cfgloader.IReloadData {
	return &ServerConfig{}
}

func (cfg *ServerConfig) ReloadCopy() {
	root.DataRoot[root.EDServerEnv] = cfg
}

// --------------------------- health check ---------------------------

// 健康检查,顺便可以做些额外的事, 比如你可以缓存一些数据到struct, 方便后续操作.
func (cfg *ServerConfig) checkValid() bool {
	cfg.SelfIP = netutil.GetLocalAddr()
	strRemoteID := fmt.Sprintf("%d", cfg.ServerID)
	remoteID, err := strconv.Atoi(strRemoteID)
	if err != nil {
		xlog.Errorf("ServerConfig calc remoteID err=%v", err)
		return false
	}
	cfg.RemoteID = int32(remoteID)
	cfg.NodeName = fmt.Sprintf("%s_%d", cfg.NodeType, cfg.RemoteID)
	root.ServerDataDir = cfg.ServerDataDir
	root.ServerConfigDir = cfg.ServerConfigDir
	return true
}

func (cfg *ServerConfig) PrintServerConfig() {
	xlog.Debugf("PrintServerConfig begin.")
	xlog.Debugf("NodeType=%s, NodeName=%s, ServerID=%d, RemoteID=%d",
		cfg.NodeType, cfg.NodeName, cfg.ServerID, cfg.RemoteID)
	xlog.Debugf("ServerDataDir=%s, ServerConfigDir=%s",
		cfg.ServerDataDir, cfg.ServerConfigDir)
	xlog.Debugf("PrintServerConfig end.")
}

// ------------- All Check is done. begin your show -----------------

// Get 服务器配置
func GetServerEnv() *ServerConfig {
	return root.DataRoot[root.EDServerEnv].(*ServerConfig)
}

func (cfg *ServerConfig) GetNodeName() string {
	return cfg.NodeName
}

func (cfg *ServerConfig) GetLogLV() int {
	return cfg.LogLevel
}
