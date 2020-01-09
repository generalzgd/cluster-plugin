/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: cluster_api.go
 * @time: 2019/12/16 5:24 下午
 * @project: testraft
 */

package plugin

import (
	`context`
	`time`

	`github.com/astaxie/beego/logs`
	`github.com/gin-gonic/gin`
	`github.com/golang/protobuf/proto`
)

type ClusterApi struct {
	clu *ClusterPlugin
}

func (p *ClusterApi) Serve(addr string) {
	//if libs.GetEnvName() != libs.ENV_ONLINE {
	//	gin.SetMode(gin.DebugMode)
	//}

	r := gin.Default()
	logs.Debug(r)
	r.GET("/agent/list", p.handleList)
	r.GET("/agent/local", p.handleLocalNode)
	r.GET("/agent/raft", p.handleRaftNode)
	r.GET("/agent/state", p.handleState)
	go r.Run(addr)
	logs.Info("start serve http:", addr)
}

func (p *ClusterApi) handleList(c *gin.Context) {
	c.JSON(200, p.clu.serfPointer.Members())
}

func (p *ClusterApi) handleLocalNode(c *gin.Context) {
	c.JSON(200, p.clu.serfPointer.LocalMember())
}

func (p *ClusterApi) handleRaftNode(c *gin.Context) {
	f := p.clu.raftPointer.GetConfiguration()
	if err := f.Error(); err != nil {
		c.String(200, "err:%v", err)
		return
	}
	c.JSON(200, f.Configuration().Servers)
	//

}

func (p *ClusterApi) handleState(c *gin.Context) {
	var list []StatusRaftReply
	for _, m := range p.clu.serfPointer.Members() {
		meta := ServerMeta{
			Name: m.Name,
			Addr: m.Addr,
		}
		if err := meta.FromMap(m.Tags); err != nil {
			continue
		}
		//
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if pool, ok := p.clu.rpcPool.Get(meta.RpcAddr()); ok {
			conn, _ := pool.Get(ctx)

			client := NewClusterRpcClient(conn.ClientConn)
			rep, err := client.SimpleCall(ctx, &CallRequest{
				Cmd: DefinedCmd_StatusRaft.String(),
			})
			if err == nil {
				obj := StatusRaftReply{}
				if err := proto.Unmarshal(rep.Data, &obj); err ==nil{
					list = append(list, obj)
				}
			}
			conn.Close()
		}
	}
	c.JSON(200, list)
}
