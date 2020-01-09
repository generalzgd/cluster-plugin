/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: rpcpool.go
 * @time: 2019/12/23 10:37 上午
 * @project: cluster-plugin
 */

package plugin

import (
	`context`
	`sync`
	`time`

	`github.com/astaxie/beego/logs`
	grpcpool `github.com/processout/grpc-go-pool`
	`google.golang.org/grpc`
)

type RpcPools struct {
	lock sync.RWMutex
	data map[string]*grpcpool.Pool
}

func NewRpcPool() *RpcPools {
	return &RpcPools{
		data: map[string]*grpcpool.Pool{},
	}
}

func (p *RpcPools) Put(addr string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	factory := func() (*grpc.ClientConn, error) {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		opts := []grpc.DialOption{
			grpc.WithReadBufferSize(10 * 1024),
			grpc.WithWriteBufferSize(10 * 1024),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(32*1024), grpc.MaxCallSendMsgSize(32*1024)),
			grpc.WithInsecure(),
		}

		conn, err := grpc.DialContext(ctx, addr, opts...)
		if err != nil {
			logs.Error("dial conn in pool. addr:%v err:%v", addr, err)
		}
		return conn, err
	}

	pool, err := grpcpool.New(factory, 1, 2, time.Second*60, time.Second*120)
	if err != nil {
		return
	}

	p.data[addr] = pool
}

func (p *RpcPools) Del(addr string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.data, addr)
}

func (p *RpcPools) Get(addr string) (*grpcpool.Pool, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	v, ok := p.data[addr]
	return v, ok
}


