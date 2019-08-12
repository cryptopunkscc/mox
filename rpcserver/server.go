package rpcserver

import (
	"log"
	"net"

	"github.com/cryptopunkscc/go-bitcoin"
	"github.com/cryptopunkscc/mox/jabber"
	"github.com/cryptopunkscc/mox/rpc"
	"github.com/cryptopunkscc/mox/services"
	"google.golang.org/grpc"
)

type RPCServer struct {
	config      *Config
	ln          bitcoin.LightningClient
	jabber      *jabber.Jabber
	MoneySender *services.MoneySender
}

func New(config *Config, j *jabber.Jabber, ln bitcoin.LightningClient, ms *services.MoneySender) *RPCServer {
	srv := &RPCServer{
		config:      config,
		ln:          ln,
		jabber:      j,
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
