package do

import (
	"io"
	"net"
)

const (
	network = "tcp"
)

// TCPProxy listen localAddr and transfer any request to remoteAddr.We can use handlers to specify one custom func to transfer data.
func TCPProxy(localAddr, remoteAddr string, handlers ...func(lconn, rconn net.Conn)) (err error) {
	laddr, err := net.ResolveTCPAddr(network, localAddr)
	if err != nil {
		return
	}
	l, err := net.ListenTCP(network, laddr)
	if err != nil {
		return
	}
	defer l.Close()

	raddr, err := net.ResolveTCPAddr(network, remoteAddr)
	if err != nil {
		return
	}

	rconn, err := net.DialTCP(network, nil, raddr)
	if err != nil {
		return
	}
	defer rconn.Close()

	for {
		var lconn net.Conn
		lconn, err = l.Accept()
		if err != nil {
			return
		}
		defer lconn.Close()

		if len(handlers) == 0 {
			TCPProxyDefaultHandler(lconn, rconn)
		} else {
			go handlers[0](lconn, rconn)
		}
	}
}

func TCPProxyDefaultHandler(lconn, rconn net.Conn) {
	go io.Copy(lconn, rconn)
	go io.Copy(rconn, lconn)
}

// TCPSend send to remote addr with handler
func TCPSend(remoteAddr string, handler func(conn net.Conn) error) (err error) {
	raddr, err := net.ResolveTCPAddr(network, remoteAddr)
	if err != nil {
		return
	}

	rconn, err := net.DialTCP(network, nil, raddr)
	if err != nil {
		return
	}
	defer rconn.Close()

	return handler(rconn)
}

// TCPRecv recv from local addr with handler, it will block on (*net.TCPListener).Accept
func TCPRecv(localAddr string, handler func(conn net.Conn)) (err error) {
	laddr, err := net.ResolveTCPAddr(network, localAddr)
	if err != nil {
		return
	}
	l, err := net.ListenTCP(network, laddr)
	if err != nil {
		return
	}
	defer l.Close()

	for {
		var lconn net.Conn
		lconn, err = l.Accept()
		if err != nil {
			return
		}
		defer lconn.Close()

		go handler(lconn)
	}
}
