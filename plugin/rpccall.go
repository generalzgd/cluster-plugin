/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: rpccall.go
 * @time: 2019/12/23 4:23 下午
 * @project: cluster-plugin
 */

package plugin

import (
	`context`
	`fmt`
	`time`

	`github.com/golang/protobuf/proto`
	`github.com/golang/protobuf/ptypes`
)

func (p *ClusterPlugin) RpcCall(addr, cmd string, id uint32, msg proto.Message) (*CallReply, error) {
	if pool, ok := p.rpcPool.Get(addr); ok {
		// buf, err := proto.Marshal(msg)
		an, err := ptypes.MarshalAny(msg)
		if err == nil {
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			conn, _ := pool.Get(ctx)
			defer conn.Close()
			client := NewClusterRpcClient(conn.ClientConn)

			req := &CallRequest{
				Cmd:  cmd,
				Id:   id,
				Data: an,
			}
			return client.SimpleCall(ctx, req)
		}
	}
	return &CallReply{}, fmt.Errorf("cannot find rpc address")
}

func (p *ClusterPlugin) leaderSwift(svrId string) {
	p.CurrentLeaderId = svrId
	p.callback.OnLeaderSwift(svrId)
	/*for _, m := range p.serfPointer.Members() {
		meta := ServerMeta{
			Name: m.Name,
			Addr: m.Addr,
		}
		if err := meta.FromMap(m.Tags); err != nil {
			continue
		}

		obj := &LeaderSwiftRequest{
			SvrID: svrId,
		}

		_, err := p.RpcCall(meta.RpcAddr(), DefinedCmd_LeaderSwift.String(), 0, obj)
		if err != nil {
			continue
		}
	}*/
}
