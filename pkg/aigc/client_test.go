package aigc

import (
	"fmt"
	"testing"

	"friedbot/pkg/config"
)

func TestGetCompletionsChat(t *testing.T) {
	err := config.InitConfig()
	if err != nil {
		t.Errorf("InitConfig() error = %v", err)
		return
	}
	err = InitClient()
	if err != nil {
		t.Errorf("InitClient() error = %v", err)
		return
	}
	type args struct {
		req *Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				req: &Request{
					Messages: []Message{
						NewSystemMessage("你是一个测试AI, 需要配合测试", "system"),
						NewUserMessage("你好, 这里是一条测试消息, 希望得到你顺利的回复", "user"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCompletionChat(tt.args.req)
			if err != nil {
				t.Errorf("GetCompletionChat() error = %v", err)
				return
			} else {
				t.Logf("GetCompletionChat() got = %v", got)
			}
		})
	}
}

func TestStreamChat(t *testing.T) {
	err := config.InitConfig()
	if err != nil {
		t.Errorf("InitConfig() error = %v", err)
		return
	}
	err = InitClient()
	if err != nil {
		t.Errorf("InitClient() error = %v", err)
		return
	}
	type args struct {
		req *Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				req: &Request{
					Messages: []Message{
						NewSystemMessage("你是一个测试AI, 需要配合测试", "system"),
						NewUserMessage("你好, 这里是一条测试流式传输的消息, 请你收到消息后回复阿拉伯数字1-10, 按空格分隔", "user"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := GetStreamChat(tt.args.req)
			if err != nil {
				t.Errorf("GetCompletionChat() error = %v", err)
				return
			} else {
				t.Logf("GetCompletionChat() got = %v", stream)
				RangeStream(stream, func(chunk string) bool {
					fmt.Println("---")
					fmt.Println(chunk)
					return true
				})
			}
		})
	}
}
