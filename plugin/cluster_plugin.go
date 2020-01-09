/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: cluster.go
 * @time: 2019/12/11 5:51 下午
 * @project: testraft
 */

package plugin

import (
	`context`
	`crypto/md5`
	`fmt`
	`io`
	`log`
	`os`
	`path/filepath`
	`strconv`
	`strings`
	`time`

	`github.com/astaxie/beego/logs`
	libs `github.com/generalzgd/comm-libs`
	`github.com/golang/protobuf/proto`
	`github.com/hashicorp/memberlist`
	`github.com/hashicorp/raft`
	raftboltdb `github.com/hashicorp/raft-boltdb`
	`github.com/hashicorp/serf/serf`

	`github.com/generalzgd/cluster-plugin/logadapter`
)

type ClusterPlugin struct {
	args     ClusterArgs
	serfAddr string
	id       string
	raftAddr string

	serfConfig  *serf.Config
	serfPointer *serf.Serf
	raftPointer *raft.Raft
	raftDb      *raftboltdb.BoltStore
	httpAddr    string
	//
	serfEvents chan serf.Event
	obsNotify  chan raft.Observation
	// 定义所有的节点都是server节点，后期考虑加入client节点
	isLeader     bool
	api          *ClusterApi
	currentState raft.RaftState
	started      bool
	firstChange  bool
	//
	callback ClusterCallback
	rpcPool  *RpcPools
	// 当前集群的服务ID。raft.serverID
	CurrentLeaderId   string
	CurrentServerMeta ServerMeta
	rpcHandlers       map[string]RpcHandler
	//
	logWriters []io.WriteCloser
}

func makeMyLogger(loggerName string) io.WriteCloser {
	dir, _ := os.Getwd()
	filename := filepath.Join(dir, "logs", fmt.Sprintf("%s.log", loggerName))
	// logs.Info("make file logger:", filename)
	logger := &logadapter.FileWriter{}
	logger.Init(filename)
	return logger // log.New(logger, "", log.LstdFlags|log.Lshortfile)
}

/*func makeRaftLogger(loggerName string) io.Writer {
	filename := fmt.Sprintf("logs/%s.log", loggerName)
	logger := &logadapter.FileWriter{}
	logger.Init(filename)
	return logger
}*/

func CreateCluster(args ClusterArgs, callback ClusterCallback) (*ClusterPlugin, error) {
	if err := args.validate(); err != nil {
		return nil, err
	}
	if callback == nil {
		callback = &emptyCallback{}
	}

	serfEvents := make(chan serf.Event, 1)
	serfAddr := args.getSerfAddr()

	memberLogger := makeMyLogger("members")
	memberListConfig := memberlist.DefaultLANConfig()
	memberListConfig.BindAddr = args.LocIp
	memberListConfig.BindPort = args.SerfPort
	// memberListConfig.LogOutput = makeLogger("members", args.GetLogLevel()) // logs.GetBeeLogger()
	memberListConfig.Logger = log.New(memberLogger, "", log.LstdFlags|log.Lshortfile)

	serfLogger := makeMyLogger("serf")
	serfConfig := serf.DefaultConfig()
	serfConfig.NodeName = serfAddr
	serfConfig.EventCh = serfEvents
	serfConfig.MemberlistConfig = memberListConfig
	// serfConfig.LogOutput = makeLogger("serf", args.GetLogLevel()) // logs.GetBeeLogger()
	serfConfig.Logger = log.New(serfLogger, "", log.LstdFlags|log.Lshortfile)

	meta := ServerMeta{
		Role:      args.Role,
		SerfPort:  args.SerfPort,
		RaftPort:  args.RaftPort,
		HttpPort:  args.HttpPort,
		RpcPort:   args.RpcPort,
		Bootstrap: args.Bootstrap,
		Except:    args.Except,
		NonVoter:  false,
	}
	if args.Agent == NodeTypeClient {
		meta.NonVoter = true
	}
	serfConfig.Tags = meta.ToMap()

	s, err := serf.Create(serfConfig)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	//
	workDir, err := os.Getwd()
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	raftAddr := args.getRaftAddr()
	id := fmt.Sprintf("%x", md5.Sum([]byte(raftAddr)))
	//
	dataDir := filepath.Join(workDir, id)
	if err := os.RemoveAll(dataDir + "/"); err != nil {
		logs.Error(err)
		return nil, err
	}
	if err := os.MkdirAll(dataDir, 0777); err != nil {
		logs.Error(err)
		return nil, err
	}
	//
	raftDBPath := filepath.Join(dataDir, "raft.db")
	raftDB, err := raftboltdb.NewBoltStore(raftDBPath)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(dataDir, 1, logs.GetBeeLogger())
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	trans, err := raft.NewTCPTransport(raftAddr, nil, 3, 10*time.Second, logs.GetBeeLogger())
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	raftLogger := makeMyLogger("raft")
	c := raft.DefaultConfig()
	c.LogOutput = raftLogger // logs.GetBeeLogger()
	// c.Logger = makeRaftLogger("raft") //
	// c.LogLevel = "error"
	c.LocalID = raft.ServerID(raftAddr)
	r, err := raft.NewRaft(c, &fsm{}, raftDB, raftDB, snapshotStore, trans)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	obsNotify := make(chan raft.Observation, 16)
	r.RegisterObserver(raft.NewObserver(obsNotify,
		false,
		func(o *raft.Observation) bool {
			return true
		}),
	)
	//
	cluster := &ClusterPlugin{
		args:              args,
		serfConfig:        serfConfig,
		serfPointer:       s,
		serfAddr:          serfAddr,
		id:                libs.Md5(raftAddr),
		raftAddr:          raftAddr,
		raftPointer:       r,
		raftDb:            raftDB,
		httpAddr:          args.getHttpAddr(),
		serfEvents:        serfEvents,
		obsNotify:         obsNotify,
		callback:          callback,
		rpcPool:           NewRpcPool(),
		firstChange:       true,
		CurrentServerMeta: meta,
		logWriters:        []io.WriteCloser{memberLogger, serfLogger, raftLogger},
	}
	cluster.initRpcHandlers()

	// 没有期望启动数量
	if args.Agent == NodeTypeServer && args.Bootstrap {
		n, err := s.Join(args.Peers, false)
		if err != nil {
			logs.Error("serf member join num:%d, error:%v", n, err)
			return nil, err
		}
		// 直接启动
		bootstrapConfig := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(raftAddr),
					Address:  raft.ServerAddress(raftAddr),
				},
			}}
		f := r.BootstrapCluster(bootstrapConfig)
		if err := f.Error(); err != nil {
			logs.Error(err)
			return nil, err
		}
		cluster.started = true
	} else {
		go cluster.waitSerfJoin(args.Peers)
	}

	cluster.api = &ClusterApi{cluster}
	cluster.api.Serve(cluster.httpAddr)
	cluster.startRpcServer()
	//
	go cluster.readSerfEvent(serfEvents, r)
	// go cluster.readRaftNotify(r.LeaderCh(), r)
	go cluster.readObserver(obsNotify)

	return cluster, nil
}

func (p *ClusterPlugin) Shutdown() {
	p.raftPointer.Shutdown()
	p.serfPointer.Shutdown()
	close(p.serfEvents)
	close(p.obsNotify)

	// 关闭日志文件
	for _, it := range p.logWriters {
		it.Close()
	}
}

func (p *ClusterPlugin) MemberList() []serf.Member {
	return p.serfPointer.Members()
}

func (p *ClusterPlugin) Serf() *serf.Serf {
	return p.serfPointer
}

func (p *ClusterPlugin) Raft() *raft.Raft {
	return p.raftPointer
}

// 广播自定义消息
func (p *ClusterPlugin) BroadcastCustomMsg(msg []byte) error {
	return p.serfPointer.UserEvent("CustomActionMsg", msg, true)
}

// 执行交互消息，设置参数可以过滤节点、超时、ack等
func (p *ClusterPlugin) Query(name string, msg []byte, param *serf.QueryParam) (*serf.QueryResponse, error) {
	return p.serfPointer.Query(name, msg, param)
}

func (p *ClusterPlugin) waitSerfJoin(peers []string) {
	for {
		n, err := p.serfPointer.Join(peers, false)
		if err == nil {
			logs.Info("serf member join num:%d", n)
			break
		} else {
			logs.Error("serf member join num:%d, error:%v", n, err)
			// 如果某个节点加入失败，则等待一会重新请求加入
			time.Sleep(time.Second * 3)
		}
	}
}

func (p *ClusterPlugin) readObserver(ch chan raft.Observation) {
	for ob := range ch {
		// r := ob.Raft
		d := ob.Data
		switch dd := d.(type) {
		case *raft.RequestVoteRequest:
			logs.Info("RequestVoteRequest:", dd)
		case raft.RaftState:
			// logs.Info("RaftState:", dd)
			p.onStateChange(dd)
		case raft.PeerObservation:
			logs.Info("PeerObservation:", dd)
			// logs.Info("logs state:", p.raftPointer.String())
			p.onStateChange(p.raftPointer.State())
		case raft.LeaderObservation:
			logs.Info("LeaderObservation:", dd, fmt.Sprintf("%v", d))
			// logs.Info("logs state:", p.raftPointer.String())
			p.onStateChange(p.raftPointer.State())
		}
	}
}

/*func (p *ClusterPlugin) readRaftNotify(ch <-chan bool, r *raft.Raft) {
	for stat := range ch {
		p.isLeader = stat
		logs.Info("raft leader stat:", stat)
		if stat {
			p.callback.OnImLeader()
		} else {
			p.callback.OnImFollower()
		}
	}
}*/

func (p *ClusterPlugin) onStateChange(stat raft.RaftState) {
	logs.Info("onStateChange: %v => %v, firstChange: %v", p.currentState, stat, p.firstChange)
	if p.currentState == stat && !p.firstChange {
		return
	}
	p.firstChange = false

	p.currentState = stat
	logs.Info("onStateChange: %v", p.currentState)
	if stat == raft.Leader {
		p.isLeader = true
	} else {
		p.isLeader = false
	}
	switch p.currentState {
	case raft.Candidate:
		p.callback.OnImCandidate()
	case raft.Follower:
		p.callback.OnImFollower()
	case raft.Leader:
		p.callback.OnImLeader()
		//
		p.leaderSwift(p.raftAddr)
	}
}

func (p *ClusterPlugin) readSerfEvent(ch chan serf.Event, r *raft.Raft) {
	for ev := range ch {
		switch dev := ev.(type) {
		case serf.MemberEvent:
			p.onMemberEvent(dev)
		case serf.UserEvent:
			p.onUserEvent(dev)
		case *serf.Query:
			p.onQueryEvent(dev)
		}
	}
}

func (p *ClusterPlugin) onMemberEvent(ev serf.MemberEvent) {
	leader := p.raftPointer.VerifyLeader()
	p.isLeader = leader.Error() == nil

	callbacks := map[serf.EventType]func([]serf.Member){
		serf.EventMemberJoin:   p.handleSerfJoin,
		serf.EventMemberLeave:  p.handleSerfLeave,
		serf.EventMemberFailed: p.handleSerfLeave,
		serf.EventMemberReap:   p.handleSerfLeave,
	}

	if f, ok := callbacks[ev.EventType()]; ok {
		go f(ev.Members)
	}
}

func (p *ClusterPlugin) maybeBootstrap() {
	index, err := p.raftDb.LastIndex()
	if err != nil {
		logs.Error("cluster: Failed to read last raft index: %v", err)
		return
	}
	if index != 0 {
		logs.Info("cluster: Raft data found, disabling bootstrap mode")
		p.args.Except = 0
		return
	}

	// Scan for all the known servers.
	members := p.serfPointer.Members()
	var servers []ServerMeta
	voters := 0
	for _, member := range members {
		it := ServerMeta{}
		it.Name = member.Name
		it.Addr = member.Addr
		if err := it.FromMap(member.Tags); err != nil {
			continue
		}
		if it.Except != 0 && it.Except != p.args.Except {
			logs.Error("cluster: Member %v has a conflicting expect value. All nodes should expect the same number.", member)
			return
		}
		if it.Bootstrap {
			logs.Error("cluster: Member %v has bootstrap mode. Expect disabled.", member)
			return
		}
		if !it.NonVoter {
			voters++
		}
		servers = append(servers, it)
	}

	// Skip if we haven't met the minimum expect count.
	if voters < p.args.Except {
		return
	}

	// Query each of the servers and make sure they report no Raft peers.
	for _, server := range servers {
		var peers []string

		// Retry with exponential backoff to get peer status from this server
		for attempt := 0; attempt < p.args.MaxPeerRetries; attempt++ {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if pool, ok := p.rpcPool.Get(server.RpcAddr()); ok {
				conn, _ := pool.Get(ctx)

				client := NewClusterRpcClient(conn.ClientConn)
				rep, err := client.SimpleCall(ctx, &CallRequest{
					Cmd: DefinedCmd_StatusPeer.String(),
				})
				if err == nil {
					obj := StatusPeerReply{}
					if err := proto.Unmarshal(rep.Data, &obj); err == nil {
						peers = obj.Peers
					}
				}
				conn.Close()
				if err == nil {
					break
				}
			}

			nextRetry := (1 << uint(attempt)) * p.args.PeerRetryBase
			time.Sleep(nextRetry)
		}

		if len(peers) > 0 {
			logs.Info("cluster: Existing Raft peers reported by %s, disabling bootstrap mode", server.Name)
			p.args.Except = 0
			return
		}
	}

	// Attempt a live bootstrap!
	var configuration raft.Configuration
	var addrs []string

	for _, server := range servers {
		addr := server.Addr.String() + ":" + strconv.Itoa(server.RaftPort)
		addrs = append(addrs, addr)
		suffrage := raft.Voter
		if server.NonVoter {
			suffrage = raft.Nonvoter
		}
		peer := raft.Server{
			ID:       raft.ServerID(addr),
			Address:  raft.ServerAddress(addr),
			Suffrage: suffrage,
		}
		configuration.Servers = append(configuration.Servers, peer)
	}
	logs.Info("cluster: Found expected number of peers, attempting bootstrap: %s",
		strings.Join(addrs, ","))
	future := p.raftPointer.BootstrapCluster(configuration)
	if err := future.Error(); err != nil {
		logs.Info("cluster: Failed to bootstrap cluster: %v", err)
	}

	p.args.Except = 0
	p.started = true
}

func (p *ClusterPlugin) handleSerfJoin(members []serf.Member) {
	logs.Debug("handleSerfJoin:", members)

	for _, m := range members {
		meta := ServerMeta{
			Name: m.Name,
			Addr: m.Addr,
		}
		if err := meta.FromMap(m.Tags); err != nil {
			continue
		}
		logs.Info("cluster: Adding LAN server %v", meta)

		p.rpcPool.Put(meta.RpcAddr())

		// If we're still expecting to bootstrap, may need to handle this.
		if p.args.Except != 0 {
			p.maybeBootstrap()
		} else {
			if p.isLeader {
				logs.Info("cluster: Leader add non or voter to raft")
				if meta.NonVoter {
					p.raftPointer.AddNonvoter(raft.ServerID(meta.RaftAddr()), raft.ServerAddress(meta.RaftAddr()), 0, 0)
				} else {
					p.raftPointer.AddVoter(raft.ServerID(meta.RaftAddr()), raft.ServerAddress(meta.RaftAddr()), 0, 0)
				}
			}
		}
	}
}

func (p *ClusterPlugin) handleSerfLeave(members []serf.Member) {
	for _, member := range members {
		meta := ServerMeta{
			Name: member.Name,
			Addr: member.Addr,
		}
		if err := meta.FromMap(member.Tags); err != nil {
			continue
		}
		p.rpcPool.Del(meta.RpcAddr())

		serfPeer := meta.SerfAddr() // fmt.Sprintf("%s:%d", member.Addr.String(), member.Port)
		//
		if err := p.serfPointer.RemoveFailedNode(serfPeer); err != nil {
			logs.Error("remove failed node:", err)
		}
		p.callback.OnNodeLeave(serfPeer)

		raftAddr := meta.RaftAddr()
		if p.isLeader && raftAddr != "" {
			if f := p.raftPointer.RemoveServer(raft.ServerID(raftAddr), 0, 0); f.Error() != nil {
				logs.Error("RemoveServer: ", f.Error())
				if f := p.raftPointer.RemovePeer(raft.ServerAddress(raftAddr)); f.Error() != nil {
					logs.Error("RemovePeer:", f.Error())
				}
			}
		} else {
			members := p.serfPointer.Members()
			aliveMembers := make([]serf.Member, 0, len(members))
			for _, m := range members {
				if m.Status == serf.StatusAlive {
					aliveMembers = append(aliveMembers, m)
				}
			}
			if len(aliveMembers) == 1 {
				addr := fmt.Sprintf("%s:%d", aliveMembers[0].Addr.String(), aliveMembers[0].Port)
				if addr == p.serfAddr {
					logs.Error("leader is gone... may be split brain")
					//
					p.callback.OnWarn(serf.UserEvent{
						Name:    "CustomActionWarn",
						Payload: []byte("leader is gone, may be brain split"),
					})
				}
			}
		}
	}
}

func (p *ClusterPlugin) onUserEvent(ev serf.UserEvent) {
	logs.Info(ev)
	handles := map[string]func([]byte){
		"CustomActionMsg": p.onUserEventCustomMsg,
	}
	f, ok := handles[ev.Name]
	if !ok {
		logs.Error("cannot find target handler by %s", ev.Name)
		return
	}
	f(ev.Payload)
}

func (p *ClusterPlugin) onUserEventCustomMsg(payload []byte) {
	p.callback.OnCustomMsg(payload)
}

// 处理query
func (p *ClusterPlugin) onQueryEvent(ev *serf.Query) {
	handles := map[string]func(*serf.Query){
	}
	f, ok := handles[ev.Name]
	if ok {
		f(ev)
		return
	}
	p.callback.OnQuery(ev)
}

/*func (p *ClusterPlugin) onGetStatusPeer(ev *serf.Query) {
	var peers []string

	f := p.raftPointer.GetConfiguration()
	if err := f.Error(); err == nil {
		for _, s := range f.Configuration().Servers {
			peers = append(peers, string(s.Address))
		}
	}
	if len(peers) > 0 {
		bts, err := json.Marshal(peers)
		if err == nil {
			ev.Respond(bts)
			return
		}
	}
	ev.Respond([]byte("[]"))
}*/
