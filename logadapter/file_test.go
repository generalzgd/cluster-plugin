/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: file_test.go.go
 * @time: 2020/1/19 10:45 上午
 * @project: cluster-plugin
 */

package logadapter

import (
	`fmt`
	`io`
	`os`
	`path/filepath`
	`testing`
)

var (
	logger *FileWriter
)

func makeMyLogger(loggerName string) io.WriteCloser {
	dir, _ := os.Getwd()
	filename := filepath.Join(dir, fmt.Sprintf("%s.log", loggerName))
	// logs.Info("make file logger:", filename)
	logger = &FileWriter{}
	logger.Init(filename)
	return logger // log.New(logger, "", log.LstdFlags|log.Lshortfile)
}

func TestFileWriter_Write(t *testing.T) {

	makeMyLogger("test")

	type args struct {
		msg []byte
	}
	tests := []struct {
		name    string
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:"TestFileWriter_Write",
			args:args{msg:[]byte("new line message")},
			wantN: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotN, err := logger.Write(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Write() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}