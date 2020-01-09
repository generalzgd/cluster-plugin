/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: def.go
 * @time: 2019/12/19 5:03 下午
 * @project: cluster-plugin
 */

package plugin

import (
	`errors`
	`fmt`
	`net`
	`strconv`
	`time`

	`github.com/astaxie/beego/logs`
	libs `github.com/generalzgd/comm-libs`
	`github.com/generalzgd/comm-libs/slice`
)

var (
	SamePortErr  = errors.New("some port is same")
	PeerAddrErr  = errors.New("some address of peer are error")
	PeerEmptyErr = errors.New("when except not 0, must has join arg")
)

// 集群节点类型（Server, Client）
type NodeType int

const (
	NodeTypeServer NodeType = 0
	NodeTypeClient NodeType = 1
)

// 集群类型（单点模式，多点模式）
type ModeType int

const (
	// 单节点
	SingleMode ModeType = 0
	// 多节点集群
	MultiMode ModeType = 1
)

type ServerMeta struct {
	Name      string
	Role      string
	Addr      net.IP
	SerfPort  int
	RaftPort  int
	RpcPort   int
	HttpPort  int
	Bootstrap bool // 启动标记
	Except    int  // 期望数量
	NonVoter  bool
}

func (p *ServerMeta) SerfAddr() string {
	return fmt.Sprintf("%s:%d", p.Addr.String(), p.SerfPort)
}

func (p *ServerMeta) RaftAddr() string {
	return fmt.Sprintf("%s:%d", p.Addr.String(), p.RaftPort)
}

func (p *ServerMeta) RpcAddr() string {
	return fmt.Sprintf("%s:%d", p.Addr.String(), p.RpcPort)
}

func (p *ServerMeta) FromMap(m map[string]string) (err error) {
	if m == nil {
		return fmt.Errorf("not target server cluster")
	}
	p.Role = m["role"]

	_, bootstrapOk := m["bootstrap"]
	p.Bootstrap = bootstrapOk

	_, nonvoterOk := m["nonvoter"]
	p.NonVoter = nonvoterOk

	getIntField := func(key string) (int, error) {
		if str, ok := m[key]; ok {
			v, err := strconv.Atoi(str)
			if err != nil {
				return 0, err
			}
			return v, nil
		}
		return 0, nil
	}
	if p.Except, err = getIntField("except"); err != nil {
		return
	}
	if p.SerfPort, err = getIntField("serfport"); err != nil {
		return
	}
	if p.RaftPort, err = getIntField("raftport"); err != nil {
		return
	}
	if p.HttpPort, err = getIntField("httpport"); err != nil {
		return
	}
	if p.RpcPort, err = getIntField("rpcport"); err != nil {
		return
	}

	return nil
}

func (p *ServerMeta) ToMap() map[string]string {
	out := map[string]string{
		"role":     p.Role,
		"except":   strconv.Itoa(p.Except),
		"raftport": strconv.Itoa(p.RaftPort),
	}
	if p.Bootstrap {
		out["bootstrap"] = "1"
	}
	if p.NonVoter {
		out["nonvoter"] = "1"
	}
	if p.HttpPort > 0 {
		out["httpport"] = strconv.Itoa(p.HttpPort)
	}
	if p.SerfPort > 0 {
		out["serfport"] = strconv.Itoa(p.SerfPort)
	}
	if p.RpcPort > 0 {
		out["rpcport"] = strconv.Itoa(p.RpcPort)
	}
	return out
}

// @param agent string: 设置运行模式，client or server
// @param serfPort int: gossip协议端口，必须设置
// @param raftPort int: raft协议端口，默认为serfPort+1
// @param httpPort int: http协议端口，默认为raftPort+1
// @param except int: 期望集群的server数量。
// 					模式1:except>0,会等待所有的server节点启动才开始raft选举(都要设置peers);
// 					模式2:except=0,第一个启动便开始raft选举(不设置peers)，后面启动的节点(设置peers)需要join,
// @param peers []string: 集群节点，[]string{ip:serfPort}
type ClusterArgs struct {
	Role           string // 集群名
	LocIp          string
	BindIp         string // 默认为0.0.0.0
	Agent          NodeType
	SerfPort       int
	RaftPort       int
	RpcPort        int
	HttpPort       int
	Except         int
	Peers          []string
	Bootstrap      bool
	MaxPeerRetries int
	PeerRetryBase  time.Duration
	LogLevel       int
}

func (p *ClusterArgs) GetLogLevel() int {
	if p.LogLevel > 0 {
		return p.LogLevel
	}
	return logs.LevelInfo
}

func (p *ClusterArgs) validate() error {
	if p.SerfPort == 0 {
		return fmt.Errorf("serf(gossip) port is empty")
	}
	if p.RaftPort == 0 {
		p.RaftPort = p.SerfPort + 1
	}
	if p.RpcPort == 0 {
		p.RpcPort = p.RaftPort + 1
	}
	if p.HttpPort == 0 {
		p.HttpPort = p.RpcPort + 1
	}
	arr := []int{p.SerfPort, p.RaftPort, p.HttpPort, p.RpcPort}
	if !slice.IsEveryUniqueInt(arr) {
		return SamePortErr
	}

	p.LocIp = libs.GetInnerIp()
	if p.BindIp == "" {
		p.BindIp = p.LocIp
	}
	if p.Except == 1 {
		p.Except = 0
	}
	// 要加入的终端地址校验失败
	if !p.validatePeers(p.Peers...) {
		return PeerAddrErr
	}
	if p.Except > 0 {
		// 模式1:except>0,会等待所有的server节点启动才开始raft选举(都要设置peers);
		if len(p.Peers) <= 0 {
			return PeerEmptyErr
		}
		// 非立刻启动
		p.Bootstrap = false
	} else {
		// 非立刻启动
		if len(p.Peers) > 0 {
			// 加入已有的集群
		} else {
			p.Bootstrap = true
			p.Peers = []string{fmt.Sprintf("%s:%d", p.BindIp, p.SerfPort)}
		}
		// 模式2:except<=1,第一个启动便开始raft选举(不设置peers)，后面启动的节点(设置peers)需要join,
	}
	// p.Agent = NodeTypeServer // 定义所有的节点都是server节点，后期考虑加入client节点

	if p.MaxPeerRetries < 1 {
		p.MaxPeerRetries = 3
	}
	if p.PeerRetryBase < 1 {
		p.PeerRetryBase = time.Millisecond * 50
	}
	return nil
}

func (p *ClusterArgs) getSerfAddr() string {
	return fmt.Sprintf("%s:%d", p.BindIp, p.SerfPort)
}

func (p *ClusterArgs) getRaftAddr() string {
	return fmt.Sprintf("%s:%d", p.BindIp, p.RaftPort)
}

func (p *ClusterArgs) getHttpAddr() string {
	return fmt.Sprintf(":%d", p.HttpPort)
}

func (p *ClusterArgs) getRpcAddr() string {
	return fmt.Sprintf(":%d", p.RpcPort)
}

func (p *ClusterArgs) getModeType() ModeType {
	if p.Except > 0 {
		return MultiMode
	}
	return SingleMode
}

func (p *ClusterArgs) validatePeers(peers ...string) bool {
	for _, peer := range peers {
		_, port, err := net.SplitHostPort(peer)
		if err != nil {
			return false
		}
		if port != strconv.Itoa(p.SerfPort) {
			return false
		}
	}
	return true
}
