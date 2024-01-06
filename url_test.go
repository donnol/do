package do

import "testing"

func TestReplaceIP(t *testing.T) {
	type args struct {
		link string
		ip   string
		nip  string
	}
	tests := []struct {
		name    string
		args    args
		wantR   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				link: "http://127.0.0.1:12306",
				ip:   "127.0.0.1",
				nip:  "127.0.0.2",
			},
			wantR:   "http://127.0.0.2:12306",
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				link: "http://127.0.0.1",
				ip:   "127.0.0.1",
				nip:  "127.0.0.2",
			},
			wantR:   "http://127.0.0.2",
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				link: "http://127.0.0.1:12306/path",
				ip:   "127.0.0.1",
				nip:  "127.0.0.2",
			},
			wantR:   "http://127.0.0.2:12306/path",
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				link: "127.0.0.1",
				ip:   "127.0.0.1",
				nip:  "127.0.0.2",
			},
			wantR:   "127.0.0.2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := ReplaceIP(tt.args.link, tt.args.ip, tt.args.nip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("ReplaceIP() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
