/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: cluster_face.go
 * @time: 2019/12/18 3:04 下午
 * @project: cluster-plugin
 */

package plugin

import (
	`context`

	`github.com/astaxie/beego/logs`
	`github.com/hashicorp/serf/serf`
)

// 回调接口
type ClusterCallback interface {
	OnLeaderSwift(string)
	OnImLeader()
	OnImFollower()
	OnImCandidate()
	OnImVoter()
	OnImNonvoter()
	OnNodeJoin(string)
	OnNodeLeave(string)
	// 节点准备就绪
	OnNodeReady()
	OnWarn(serf.UserEvent)
	OnCustomMsg([]byte)
	OnQuery(*serf.Query)
	OnRpcCall(context.Context, *CallRequest) (*CallReply, error)
}

//
type emptyCallback struct {
}

func (p *emptyCallback) OnLeaderSwift(id string) {
	logs.Debug("OnLeaderSwift:", id)
}

func (p *emptyCallback) OnImLeader() {
	logs.Debug("OnImLeader")
}

func (p *emptyCallback) OnImFollower() {
	logs.Debug("OnImFollower")
}

func (p *emptyCallback) OnImCandidate() {
	logs.Debug("OnImCandidate")
}

func (p *emptyCallback) OnNodeJoin(peer string) {
	logs.Debug("OnNodeJoin:", peer)
}

func (p *emptyCallback) OnNodeLeave(peer string) {
	logs.Debug("OnNodeLeave:", peer)
}

func (p *emptyCallback) OnWarn(ev serf.UserEvent) {
	logs.Debug("OnWarn:", ev)
}

func (p *emptyCallback) OnCustomMsg(bts []byte) {
	logs.Debug("OnCustomMsg:", string(bts))
}

func (p *emptyCallback) OnQuery(q *serf.Query) {
	logs.Debug("OnQuery:", q)
}

func (p *emptyCallback) OnImVoter() {
	logs.Debug("OnImVoter")
}

func (p *emptyCallback) OnImNonvoter() {
	logs.Debug("OnImNonvoter")
}

func (p *emptyCallback) OnRpcCall(ctx context.Context, req *CallRequest) (*CallReply, error) {
	return &CallReply{Cmd: req.Cmd, Id: req.Id}, nil
}

func (p *emptyCallback) OnNodeReady() {

}