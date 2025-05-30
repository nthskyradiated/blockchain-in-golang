package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/nthskyradiated/blockchain-in-golang/utils"
	// "crypto/sha256"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce	int
}

// func (b *Block) DeriveHash() {
// 	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
// 	hash := sha256.Sum256(info)
// 	b.Hash = hash[:]
// }

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func GenesisBlock() *Block {
	return CreateBlock("Genesis block", []byte{})
}

func (b *Block) Serialize() []byte {
	var encoded bytes.Buffer
	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(b)
	utils.HandleError(err)
	return encoded.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.HandleError(err)
	return &block
}