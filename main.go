/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: main.go
 * @time: 2019/12/10 3:40 下午
 * @project: testraft
 */

package main

import (
	`flag`
	`os`
	`os/signal`
	`strings`
	`syscall`

	`github.com/astaxie/beego/logs`

	`github.com/generalzgd/cluster-plugin/plugin`
)

var (
	join     string
	serfPort int
	raftPort int
	httpPort int
	rpcPort  int
	agent    string
	except   int
)

func init() {
	flag.StringVar(&join, "join", "", "127.0.0.1:1025,127.0.0.1:1035, gossip协议地址的组合")
	flag.IntVar(&serfPort, "serf", 7770, "gossip协议端口")
	flag.IntVar(&raftPort, "raft", 0, "raft协议端口，可选，默认为sert+1")
	flag.IntVar(&rpcPort, "rpc", 0, "rpc协议端口，可选，默认为raft+1")
	flag.IntVar(&httpPort, "http", 0, "http协议端口，可选，默认为rpc+1")
	flag.IntVar(&except, "except", 0, "except server node num, recommend is 3 or 5")
	flag.StringVar(&agent, "agent", "server", "传递server,或client")

	logger := logs.GetBeeLogger()
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(3)
}

func main() {
	flag.Parse()

	var peers []string
	if len(join) > 0 {
		peers = strings.Split(join, ",")
	}

	agentType := plugin.NodeTypeServer
	if agent == "client" {
		agentType = plugin.NodeTypeClient
	}

	args := plugin.ClusterArgs{
		Role:     "cluster-plugin",
		Agent:    agentType,
		SerfPort: serfPort,
		RaftPort: raftPort,
		RpcPort:  rpcPort,
		HttpPort: httpPort,
		Except:   except,
		Peers:    peers,
	}

	cluster, err := plugin.CreateCluster(args, nil)
	if err != nil {
		return
	}
	//
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)
	sig := <-chSig
	logs.Info("exit with:", sig)
	//
	cluster.Shutdown()
}
