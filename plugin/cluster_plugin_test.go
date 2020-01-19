/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: cluster_plugin_test.go.go
 * @time: 2019/12/16 5:28 下午
 * @project: testraft
 */

package plugin

import (
	`encoding/json`
	`testing`
)

func TestJsonMarshal(t *testing.T) {
	var peer []string
	bts, err := json.Marshal(peer)
	t.Logf("%s, %v", string(bts), err)
}

func TestCreateCluster(t *testing.T) {
	type args struct {
		createArgs ClusterArgs
		callback   ClusterCallback
	}
	tests := []struct {
		name    string
		args    args
		want    *ClusterPlugin
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "CreateCluster",
			args: args{
				createArgs: ClusterArgs{
					SerfPort: 7777,
					RaftPort: 7778,
					RpcPort:  7779,
					HttpPort: 7780,
					Agent:    0,
					Peers:    nil,
				},
				callback: &emptyCallback{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateCluster(tt.args.createArgs, tt.args.callback, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
