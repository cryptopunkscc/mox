package rpcserver

import (
	"context"

	"github.com/cryptopunkscc/mox/rpc"
)

func (srv *RPCServer) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	srv.MoneySender.Send(req.Address, int(req.Amount))

	return &rpc.SendResponse{}, nil
}
