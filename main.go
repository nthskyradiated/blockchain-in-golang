package main

import (
	"fmt"
	"strconv"

	"github.com/nthskyradiated/blockchain-in-golang/blockchain"
)

func main() {
	bc := blockchain.NewBlockChain()

	bc.AddBlock("First Block after Genesis")
	bc.AddBlock("Second Block after Genesis")
	bc.AddBlock("Third Block after Genesis")

	for _, bc := range bc.Blocks {
		fmt.Printf("Previous Hash: %x\n", bc.PrevHash)
		fmt.Printf("Data in Block: %s\n", bc.Data)
		fmt.Printf("Hash: %x\n", bc.Hash)

		pow := blockchain.NewProofOfWork(bc)
			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			fmt.Println()
	}

}