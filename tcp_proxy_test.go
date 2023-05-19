package do

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

func TestTCPProxy(t *testing.T) {
	type args struct {
		localAddr  string
		remoteAddr string
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
				localAddr:  ":32288",
				remoteAddr: ":32299",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// start remote server
			go startRemoteServer(tt.args.remoteAddr)
			time.Sleep(time.Millisecond * 100)

			go func() {
				if err := TCPProxy(tt.args.localAddr, tt.args.remoteAddr); (err != nil) != tt.wantErr {
					t.Errorf("TCPProxy() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			time.Sleep(time.Millisecond * 100)

			// access local addr, and it will transfer data to remote
			if err := localClient(tt.args.localAddr); err != nil {
				t.Errorf("send request to local addr failed: %v", err)
			}
		})
	}
}

func TestTCPProxySpecifyHandler(t *testing.T) {
	type args struct {
		localAddr  string
		remoteAddr string
		handler    func(lconn net.Conn, rconn net.Conn)
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
				localAddr:  ":32388",
				remoteAddr: ":32399",
				handler: func(lconn, rconn net.Conn) {
					// log.Printf("custom handler begin\n")
					var i int

					TCPProxyDefaultHandler(lconn, rconn)

					i++
					// log.Printf("custom handler end\n")
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// start remote server
			go startRemoteServer(tt.args.remoteAddr)
			time.Sleep(time.Millisecond * 100)

			go func() {
				if err := TCPProxy(tt.args.localAddr, tt.args.remoteAddr, tt.args.handler); (err != nil) != tt.wantErr {
					t.Errorf("TCPProxy() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			time.Sleep(time.Millisecond * 100)

			// access local addr, and it will transfer data to remote
			if err := localClient(tt.args.localAddr); err != nil {
				t.Errorf("send request to local addr failed: %v", err)
			}
		})
	}
}

// NOTE: just test the connection in proxy, so just use simple data for read and write which doesn't handle the longer data than 1024.

var (
	send = []byte("hello")
	recv = []byte("ack")
)

func startRemoteServer(localAddr string) (err error) {
	TCPRecv(localAddr, func(conn net.Conn) {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("server read failed: %v\n", err)
			}
			if !bytes.Equal(buf[:n], send) {
				panic(fmt.Errorf("server read wrong data: %q != %q", buf, send))
			}
			// log.Printf("server read %d bytes: %s\n", n, buf)

			{
				buf := recv
				n, err := conn.Write(buf)
				if err != nil {
					log.Printf("server write failed: %v\n", err)
				}
				_ = n
				// log.Printf("server write %d bytes: %s\n", n, buf)
			}
		}
	})
	return nil
}

func localClient(remoteAddr string) (err error) {
	return TCPSend(
		remoteAddr,
		func(conn net.Conn) error {
			for i := 0; i < 5; i++ {
				buf := send
				n, err := conn.Write(buf)
				if err != nil {
					return err
				}
				_ = n
				// log.Printf("client write %d bytes: %s\n", n, buf)

				{
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						return err
					}
					if !bytes.Equal(buf[:n], recv) {
						panic(fmt.Errorf("client read wrong data: %q != %q", buf, recv))
					}
					// log.Printf("client read %d bytes: %s\n", n, buf)
					// log.Println("=== ===")
				}

				time.Sleep(time.Millisecond * 500)
			}

			return nil
		},
	)
}
