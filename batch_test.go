package do

import (
	"fmt"
	"testing"
)

func TestBatchRun(t *testing.T) {
	type args struct {
		s        []int
		batchNum int
		handler  func([]int) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				s:        []int{1, 2, 3},
				batchNum: 1,
				handler: func(part []int) error {
					if len(part) == 0 {
						return fmt.Errorf("len is 0")
					}
					// fmt.Println(part)
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BatchRun(tt.args.s, tt.args.batchNum, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("BatchRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStreamRun(t *testing.T) {
	type args struct {
		s        chan int
		batchNum int
		handler  func([]int) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				s:        make(chan int, 5),
				batchNum: 1,
				handler: func(part []int) error {
					if len(part) == 0 {
						return fmt.Errorf("len is 0")
					}
					// fmt.Println(s)
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				s:        make(chan int, 5),
				batchNum: 2,
				handler: func(part []int) error {
					if len(part) == 0 {
						return fmt.Errorf("len is 0")
					}
					// fmt.Println(s)
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				for i := 0; i < 10; i++ {
					tt.args.s <- i
				}
				close(tt.args.s)
			}()
			if err := StreamRun(tt.args.s, tt.args.batchNum, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("StreamRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
