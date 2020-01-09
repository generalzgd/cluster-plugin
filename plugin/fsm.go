/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: fsm.go
 * @time: 2019/12/19 5:03 下午
 * @project: cluster-plugin
 */

package plugin

import (
	`io`

	`github.com/astaxie/beego/logs`
	`github.com/hashicorp/raft`
)

type fsm struct {
}

func (f *fsm) Apply(log *raft.Log) interface{} {
	logs.Debug("Apply")
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	logs.Debug("Snapshot")
	return nil, nil
}

func (f *fsm) Restore(io.ReadCloser) error {
	logs.Debug("Restore")
	return nil
}
