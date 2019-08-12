package main

import (
	"flag"
	"fmt"
	"os"

	pb "github.com/cryptopunkscc/mox/rpc"
	"google.golang.org/grpc"
)

type Command func(rpc pb.WalletClient, args []string) error

var address string

func init() {
	flag.StringVar(&address, "s", "127.0.0.1:50000", "RPC server to connect to")
}

func connect() pb.WalletClient {
	grpc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return pb.NewWalletClient(grpc)
}

func main() {
	flag.Parse()

	handlers := make(map[string]Command)
	handlers["balance"] = Balance
	handlers["send"] = Send

	rpc := connect()
	cmd := flag.Arg(0)
	if handler, ok := handlers[cmd]; ok {
		args := flag.Args()
		err := handler(rpc, args[1:len(args)])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Unknown command:", cmd)
	}

}
