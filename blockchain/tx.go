package blockchain

import (
	"bytes"
	"github.com/nthskyradiated/blockchain-in-golang/utils"
	"github.com/nthskyradiated/blockchain-in-golang/wallet"
)

type TxInput struct {
	ID       []byte
	OutIndex int
	Sig      []byte
	PubKey   []byte
}

type TxOutput struct {
	Value        int
	ScriptPubKey []byte
}

type TxOutputs struct{
	Outputs []TxOutput
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)
	return bytes.Equal(lockingHash, pubKeyHash)
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.ScriptPubKey = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	// lockingHash := wallet.PublicKeyHash(out.ScriptPubKey)
	return bytes.Equal(out.ScriptPubKey, pubKeyHash)
}

func NewTXOutput(value int, address string) *TxOutput {
	out := &TxOutput{value, nil}
	out.Lock([]byte(address))
	return out
}

func (outs TxOutputs) Serialize() []byte {
	return utils.Serialize(outs)
}
