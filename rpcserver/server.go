package rpcserver

import (
	"log"
	"net"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/mox/rpc"
	"github.com/cryptopunkscc/mox/services"
	"github.com/cryptopunkscc/mox/xmpp"
	"google.golang.org/grpc"
)

type RPCServer struct {
	config      *Config
	ln          bitcoin.LightningClient
	xmpp        *xmpp.XMPP
	MoneySender *services.MoneySender
}

func New(config *Config, j *xmpp.XMPP, ln bitcoin.LightningClient, ms *services.MoneySender) *RPCServer {
	srv := &RPCServer{
		config:      config,
		ln:          ln,
		xmpp:        j,
		MoneySender: ms,
	}
	return srv
}

func (srv *RPCServer) Run() {
	if srv.config == nil {
		log.Println("RPC server not configured")
		return
	}
	log.Println("Starting rpc server on", srv.config.bindAddress())
	tcp, err := net.Listen("tcp", srv.config.bindAddress())
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterWalletServer(grpcServer, srv)
	grpcServer.Serve(tcp)
}
