package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/nthskyradiated/blockchain-in-golang/utils"
)

type Block struct {
	Hash     []byte
	Transactions []*Transaction
	PrevHash []byte
	Nonce	int
}


func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func GenesisBlock(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) Serialize() []byte {
	return utils.Serialize(b)
}


func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.HandleError(err)
	return &block
}