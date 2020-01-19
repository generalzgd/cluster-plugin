/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: rpchandler.go
 * @time: 2019/12/23 12:01 下午
 * @project: cluster-plugin
 */

package plugin

import (
	`context`
	`net`
	`time`

	`github.com/astaxie/beego/logs`
	`github.com/golang/protobuf/ptypes`
	`google.golang.org/grpc`
	`google.golang.org/grpc/reflection`
)

type RpcHandler func(context.Context, *CallRequest) (*CallReply, error)

func (p *ClusterPlugin) initRpcHandlers() {
	p.rpcHandlers = map[string]RpcHandler{
		DefinedCmd_StatusPeer.String(): p.onStatusPeer,
		DefinedCmd_StatusRaft.String(): p.onStatusRaft,
		// DefinedCmd_LeaderSwift.String(): p.onLeaderSwift,
	}
}

func (p *ClusterPlugin) startRpcServer() {
	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(1000),
		grpc.MaxRecvMsgSize(32 * 1024),
		grpc.MaxSendMsgSize(32 * 1024),
		grpc.ReadBufferSize(8 * 1024),
		grpc.WriteBufferSize(8 * 1024),
		grpc.ConnectionTimeout(5 * time.Second),
	}
	addr := p.args.getRpcAddr()
	s := grpc.NewServer(opts...)
	RegisterClusterRpcServer(s, p)
	reflection.Register(s)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Error("failed to listen: %v", err)
	}
	go func() {
		if err := s.Serve(lis); err != nil {
			logs.Error("failed to serve: %v", err)
		}
	}()
	logs.Debug("start serve grpc.", addr)
}

func (p *ClusterPlugin) SimpleCall(ctx context.Context, req *CallRequest) (*CallReply, error) {
	// pe, _ := peer.FromContext(ctx)
	// logs.Debug("StatusPeer: req:%v, peer:%v", req, pe)
	//

	f, ok := p.rpcHandlers[req.Cmd]
	if ok {
		return f(ctx, req)
	}

	rep, err := p.callback.OnRpcCall(ctx, req)
	if rep == nil {
		rep = &CallReply{
			Cmd: req.Cmd,
			Id:  req.Id,
		}
	}
	return rep, err
}

func (p *ClusterPlugin) onStatusPeer(ctx context.Context, req *CallRequest) (rep *CallReply, err error) {
	rep = &CallReply{
		Cmd: req.Cmd,
		Id:  req.Id,
	}
	var peers []string

	f := p.raftPointer.GetConfiguration()
	if err := f.Error(); err == nil {
		for _, s := range f.Configuration().Servers {
			peers = append(peers, string(s.Address))
		}
	}
	obj := &StatusPeerReply{
		Peers: peers,
	}
	if an, err := ptypes.MarshalAny(obj); err == nil {
		rep.Data = an
	}
	return
}

func (p *ClusterPlugin) onStatusRaft(ctx context.Context, req *CallRequest) (rep *CallReply, err error) {
	rep = &CallReply{
		Cmd: req.Cmd,
		Id:  req.Id,
	}

	obj := &StatusRaftReply{
		FromRaftAddr: p.args.getRaftAddr(),
		Status:       p.raftPointer.State().String(),
	}
	if an, err := ptypes.MarshalAny(obj); err == nil {
		rep.Data = an
	}
	return
}

/*func (p *ClusterPlugin) onLeaderSwift(ctx context.Context, req *CallRequest) (rep *CallReply, err error) {
	rep = &CallReply{
		Cmd: req.Cmd,
		Id:  req.Id,
	}

	obj := LeaderSwiftRequest{}
	if err := proto.Unmarshal(req.Data, &obj); err == nil {
		p.CurrentLeaderId = obj.SvrID
		p.callback.OnLeaderSwift(obj.SvrID)
	}
	return
}*/
