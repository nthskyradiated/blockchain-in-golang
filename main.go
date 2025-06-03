package main

import (
	"os"
	"github.com/nthskyradiated/blockchain-in-golang/cli"
)


func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
