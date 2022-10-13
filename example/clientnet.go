package main

import (
	"encoding/binary"
	"fmt"
	"github.com/qixi7/xengine_core/xlog"
	"github.com/qixi7/xengine_core/xmemory"
	"github.com/qixi7/xengine_core/xmodule"
	"github.com/qixi7/xengine_core/xnet"
	"strings"
	"time"
	"xserver/example/loader"
	"xserver/example/logic"
)

type msClientHandler struct {
}

func (h *msClientHandler) OnOpen(s xnet.Sessioner, b bool) {
	xlog.InfoF("client ip=%s, sessionID=%d OnOpen", s.RemoteAddr(), s.ID())
	sliceAddr := strings.Split(s.RemoteAddr(), ":")
	onlyIP := sliceAddr[0]
	logic.GetAgentMgr().OnAgentConnect(0, s, onlyIP)
}

func (h *msClientHandler) OnClose(s xnet.Sessioner) {
	xlog.InfoF("client sessionID=%d OnClose", s.ID())
	logic.GetAgentMgr().OnAgentDisConnect(0, s.ID())
	s.Close()
}

func (h *msClientHandler) OnMessage(s xnet.Sessioner, bindata []byte, pk *xnet.PBPacket) {
	if pk != nil {
		// 修正connID. 再分发msg
		pk.ConnID = s.ID()
		logic.GetAgentMgr().DispatchMsg(0, pk)
	}
}

type gsClientBodyFormatter struct {
	pbFmt   xnet.BodyFormater
	creator xmemory.ByteSliceCreator
}

func newGsClientBodyFormatter(encoding string) *gsClientBodyFormatter {
	return &gsClientBodyFormatter{
		pbFmt: xnet.NewPBFormater(xnet.CliHeadLen, encoding),
	}
}

func (f *gsClientBodyFormatter) New() xnet.BodyFormater {
	return &gsClientBodyFormatter{
		pbFmt: f.pbFmt.New(),
	}
}

// 只有Recv协程在使用, 虽是协程环境, 但是安全
func (f *gsClientBodyFormatter) ReadBufferAlloc(len int) []byte {
	return f.creator.Create(len, len, 4096)
}

func (f *gsClientBodyFormatter) Encode(msg interface{}) ([]byte, error) {
	switch m := msg.(type) {
	case *xnet.PBPacket:
		return f.pbFmt.Encode(m)
	case *[]byte:
		return (*m)[xnet.CliHeadLen:], nil
	}
	return nil, xnet.ErrBodyEncode
}

func (f *gsClientBodyFormatter) Decode(msg []byte, s xnet.Streamer) ([]byte, *xnet.PBPacket, error, bool) {
	if len(msg) < xnet.MinCliMsgLen {
		return nil, nil, xnet.ErrPBDecodeMinSize, false
	}
	msgID := binary.BigEndian.Uint16(msg[xnet.CliHeadLen+8 : xnet.CliHeadLen+10])
	// dispatch msg
	if msgID == 1 { // echo non-proto message
		s.SendRaw(msg)
		return nil, nil, nil, false
	}
	_, pack, err, needPostMsg := f.pbFmt.Decode(msg, s)
	return nil, pack, err, needPostMsg
}

// 创建客户端TCP Server
func newClientTCPGate() xnet.Server {
	cfg := loader.GetServerEnv()
	headFormat := xnet.NewBE4ByteHeader()
	bodyFormat := newGsClientBodyFormatter("")
	timeout := time.Duration(cfg.ClientTimeout) * time.Second
	for port := cfg.PortBegin; port < cfg.PortBegin+cfg.MaxPortRange; port++ {
		s := xnet.NewTcpServer(&xnet.ServerConfig{
			Network:       "tcp",
			Addr:          fmt.Sprintf("0.0.0.0:%d", port),
			PostEvent:     true,
			Handler:       &msClientHandler{},
			Factory:       xnet.NewTLVStreamFactory(timeout, headFormat, bodyFormat),
			MaxSessionNum: uint32(cfg.MaxClient),
			QueueBufLen:   64,
			Cerificate:    nil,
		})
		if s != nil {
			// 保存正在监听的TCP地址
			cfg.NowTCPListenAddr = fmt.Sprintf("%s:%d", cfg.SelfIP, port)
			return s
		}
	}
	return nil
}

// 处理网络线程数据
type NetThread struct {
	clientTCPServer xnet.Server
}

func NewNetThread() *NetThread {
	return &NetThread{}
}

func (d *NetThread) GetTCPServer() xnet.Server {
	return d.clientTCPServer
}

func (d *NetThread) Init(selfGetter xmodule.DModuleGetter) bool {
	cliTcpServer := newClientTCPGate()
	if cliTcpServer == nil {
		xlog.Errorf("newClientTCPGate err.")
		return false
	}
	d.clientTCPServer = cliTcpServer
	return true
}

func (d *NetThread) Run(delta int64) {
	d.clientTCPServer.ProcessEvent()
}

func (d *NetThread) Destroy() {
	xlog.InfoF("client TCP channel closing")
	d.clientTCPServer.Close()
	xlog.InfoF("client channel closed")
}
