package netmodule

import (
	"encoding/binary"
	"fmt"
	"pb"
	"robot/global"
	"robot/loader"
	"time"
	"xcore/xcontainer/timer"
	"xcore/xlog"
	"xcore/xmemory"
	"xcore/xnet"
)

/*
	netmodule.go: 机器人网络模块
*/

// 网络模块
type NetModule struct {
	// gs about
	gsConn *xnet.Client // gs conn
	// bs about
	bsConn *xnet.Client // bs conn
	// common
	connected [global.EConnectTypeMax]bool
}

// --------------------------- 内置函数 ---------------------------

// --------------------------- 外置函数 ---------------------------

// init
func (c *NetModule) Init() {
	c.StartHeatBeats()
}

// 网络update
func (c *NetModule) Update() {
	if c.gsConn != nil {
		c.gsConn.ProcessEvent()
	}
	if c.bsConn != nil {
		c.bsConn.ProcessEvent()
	}
}

type robotBodyFormatter struct {
	pbFmt   xnet.BodyFormater
	creator xmemory.ByteSliceCreator
}

func newRobotBodyFormatter(headerLen int) *robotBodyFormatter {
	return &robotBodyFormatter{
		pbFmt: xnet.NewPBFormater(headerLen, ""),
	}
}

func (f *robotBodyFormatter) New() xnet.BodyFormater {
	return &robotBodyFormatter{
		pbFmt: f.pbFmt.New(),
	}
}

// 只有Recv协程在使用, 虽是协程环境, 但是安全
func (f *robotBodyFormatter) ReadBufferAlloc(len int) []byte {
	return f.creator.Create(len, len, 4096)
}

func (f *robotBodyFormatter) Encode(msg interface{}) ([]byte, error) {
	switch m := msg.(type) {
	case *xnet.PBPacket:
		return f.pbFmt.Encode(m)
	case *[]byte:
		return (*m)[xnet.CliHeadLen:], nil
	}
	return nil, xnet.ErrBodyEncode
}

func (f *robotBodyFormatter) Decode(msg []byte, s xnet.Streamer) ([]byte, *xnet.PBPacket, error, bool) {
	if len(msg) < xnet.MinCliMsgLen {
		return nil, nil, xnet.ErrPBDecodeMinSize, false
	}
	msgID := binary.BigEndian.Uint16(msg[xnet.CliHeadLen+8 : xnet.CliHeadLen+10])
	// dispatch msg
	if msgID == 1 { // echo non-proto message
		return nil, nil, nil, false
	}
	_, pack, err, needPostMsg := f.pbFmt.Decode(msg, s)
	return nil, pack, err, needPostMsg
}

// connect gs
func (c *NetModule) ConnectGS(handler xnet.SessionEventHandler) {
	// 默认gs配置
	cfg := loader.GetRobotCfg()
	headFormatter := xnet.NewBE4ByteHeader()
	msAddr := fmt.Sprintf("%s:%d", cfg.ServerIP, cfg.GsTcpPort)
	connCfg := &xnet.ClientConfig{
		Protocol:  xnet.ProtocolTCP,
		Network:   "tcp",
		Addr:      msAddr,
		PostEvent: true,
		Handler:   handler,
		Factory:   xnet.NewTLVStreamFactory(30*time.Second, headFormatter, newRobotBodyFormatter(headFormatter.HeadLen())),
		//QueueBufLen: 64,
		Timeout: 90 * time.Second,
	}
	cli, err := xnet.PostDial(connCfg)
	if err != nil {
		xlog.Errorf("Connect GS, addr=%s, err=%v", msAddr, err)
		return
	}
	c.gsConn = cli
}

// 设置连接状态
func (c *NetModule) SetConnected(connType int, connected bool) {
	c.connected[connType] = connected
}

// 获取连接状态
func (c *NetModule) GetConnected(connType int) bool {
	return c.connected[connType]
}

// 根据类型获取连接
func (c *NetModule) GetConnByType(connType int) *xnet.Client {
	switch connType {
	case global.EConnectTypeGS:
		return c.gsConn
	case global.EConnectTypeBS:
		return c.bsConn
	default:
		return nil
	}
}

// 给gs发消息
func (c *NetModule) SendGsMsg(msgID pb.Game_Msg, proto xnet.ProtoMessage) {
	newpack := xnet.New_PBPacket()
	newpack.MsgID = uint16(msgID)
	newpack.Body = proto
	c.gsConn.AsyncSend(newpack)
}

// 给bs发消息
func (c *NetModule) SendBsMsg(msgID pb.Game_Msg, proto xnet.ProtoMessage) {
	if c.bsConn == nil {
		xlog.Errorf("send msgID=%d, err, bsConn nil. connInfo=%v",
			msgID, c.connected)
		return
	}
	newpack := xnet.New_PBPacket()
	newpack.MsgID = uint16(msgID)
	newpack.Body = proto
	c.bsConn.AsyncSend(newpack)
}

// 关闭连接
func (c *NetModule) CloseConnect(connType int) {
	switch connType {
	case global.EConnectTypeGS:
		// gs
		if c.gsConn != nil {
			c.gsConn.Close()
		}
	case global.EConnectTypeBS:
		// bs
		if c.bsConn != nil {
			c.bsConn.Close()
		}
	default:
		xlog.Errorf("CloseConnect err, connType=%d, unKnown", connType)
		return
	}
}

func (c *NetModule) GetBsLocalAddr() string {
	return c.bsConn.LocalAddr()
}

// 开启心跳
func (c *NetModule) StartHeatBeats() {
	global.GetTimerMgr().Add(12*time.Second, &heartBeatsTimer{owner: c, connType: global.EConnectTypeGS})
}

// 心跳timer
type heartBeatsTimer struct {
	owner    *NetModule
	connType int
}

func (t *heartBeatsTimer) Invoke(ctl *timer.Controller, item *timer.Item) {
	if t.owner.GetConnected(t.connType) {
		newpack := xnet.New_PBPacket()
		newpack.MsgID = 1
		t.owner.GetConnByType(t.connType).AsyncSend(newpack)
	}
	ctl.Next(item)
}
