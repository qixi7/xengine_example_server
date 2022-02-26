package global

/*
	const.go: some global const define here.
*/

// 服务器帧率
const DefaultFPS = 25

// socket 连接类型
const (
	EConnectTypeGS int = iota // gs
	EConnectTypeBS            // bs
	EConnectTypeMax
)
