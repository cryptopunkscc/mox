package main

import (
	"flag"

	"github.com/cryptopunkscc/mox"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "config.json", "config file to use")
}

func main() {
	flag.Parse()
	cfg := mox.LoadConfig(configFile)
	app, err := mox.New(cfg)
	if err != nil {
		panic(err)
	}
	app.Run()
}
