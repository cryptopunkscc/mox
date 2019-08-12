package main

import (
	"context"
	"fmt"

	pb "github.com/cryptopunkscc/mox/rpc"
)

func Balance(rpc pb.WalletClient, args []string) error {
	res, err := rpc.Balance(context.Background(), &pb.BalanceRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("Channel: % 15d\n", res.LightningBalance)
	fmt.Printf("Chain:   % 15d\n", res.OnchainBalance)
	return nil
}
