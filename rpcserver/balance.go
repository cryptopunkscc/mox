package rpcserver

import (
	"context"

	"github.com/cryptopunkscc/mox/rpc"
)

func (srv *RPCServer) Balance(ctx context.Context, req *rpc.BalanceRequest) (*rpc.BalanceResponse, error) {
	return &rpc.BalanceResponse{
		LightningBalance: int64(srv.wallet.Balance()),
	}, nil
}
