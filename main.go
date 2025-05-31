package main

import (
	"os"
	"github.com/nthskyradiated/blockchain-in-golang/cli"
	// "github.com/nthskyradiated/blockchain-in-golang/wallet"
)


func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()

	// w := wallet.CreateWallet()
	// w.Address()
}
