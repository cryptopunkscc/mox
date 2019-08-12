package main

import (
	"context"
	"errors"
	"strconv"

	pb "github.com/cryptopunkscc/mox/rpc"
)

func Send(rpc pb.WalletClient, args []string) error {
	if len(args) < 2 {
		return errors.New("send <address> <amount>")
	}
	address := args[0]
	amount, _ := strconv.ParseInt(args[1], 10, 64)

	_, err := rpc.Send(context.Background(), &pb.SendRequest{
		Address: address,
		Amount:  amount,
	})
	return err
}
