package rpcserver

import (
	"github.com/cryptopunkscc/mox/wallet"
	"log"
	"net"

	"github.com/cryptopunkscc/mox/rpc"
	"github.com/cryptopunkscc/mox/xmpp"
	"google.golang.org/grpc"
)

type RPCServer struct {
	config *Config
	xmpp   *xmpp.XMPP
	wallet *wallet.Service
}

func New(config *Config, j *xmpp.XMPP, wallet *wallet.Service) *RPCServer {
	srv := &RPCServer{
		config: config,
		wallet: wallet,
		xmpp:   j,
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
