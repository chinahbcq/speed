package rpc_framework

import (
	"net"
	"os"
	"runtime"
	"strings"
)

import (
	"ublog"
)

type TcpServer struct {
	ListenAddress   string
	IdleTimeout     int
	TerminateSignal bool
	Logger          *ublog.UbLog

	Router RouterInterface
}

func (server *TcpServer) Start() {
	listener, err := net.Listen("tcp", server.ListenAddress)
	if err != nil {
		server.Logger.UbLogWarning("listen error:%v", err)
		os.Exit(1)
	}
	server.Logger.UbLogNotice("TCP: listening on %s", listener.Addr())

	for !server.TerminateSignal {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				server.Logger.UbLogWarning("temporary Accept() failure error:%s", err)
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				server.Logger.UbLogWarning("UNKOWN ERROR of Accept() error:%s", err)
			}
			continue
		}
		go server.Handle(clientConn)
	}

	server.Logger.UbLogWarning("Server: closing %s", listener.Addr())
}

func (server *TcpServer) Handle(clientConn net.Conn) {
	server.Logger.UbLogTrace("new Tcp Conn, client:%s", clientConn.RemoteAddr())
	var err error
	err = nil

	prot := NsHeadProtocol{
		IdleTimeout: server.IdleTimeout,
		Logger:      server.Logger,
		Router:      server.Router,
	}

	for err == nil {
		err = prot.IOLoop(clientConn)
	}
	clientConn.Close()
	if err != nil {
		server.Logger.UbLogWarning("[error:handle_connection] [client:%s] [error_msg:%s]", clientConn.RemoteAddr(), err)
		return
	}
}
