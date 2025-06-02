package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"runtime"
	"github.com/dgraph-io/badger"
	"github.com/nthskyradiated/blockchain-in-golang/utils"
)

const (
	dbPath = "./tmp/blocks"
	dbFile = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}


func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func (bc *BlockChain) AddBlock(transactions []*Transaction) *Block {
	var lastHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.HandleError(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
			})
		return err 
		})
	utils.HandleError(err)


	newBlock := CreateBlock(transactions, lastHash)

	err = bc.Database.Update(func(txn *badger.Txn) error {
	err := txn.Set(newBlock.Hash, newBlock.Serialize())
	utils.HandleError(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return err
	})
	utils.HandleError(err)
	return newBlock
}

func NewBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBExists() {
		fmt.Println("Blockchain already exists, loading from disk")
		runtime.Goexit()
	}
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	utils.HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
	cbtx := CoinbaseTx(address, genesisData)
		genesis := GenesisBlock(cbtx)
		fmt.Println("Genesis Block Created")
		err := txn.Set(genesis.Hash, genesis.Serialize())
		utils.HandleError(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})
	utils.HandleError(err)
	bc := BlockChain{lastHash, db}
	return &bc
}

func ContinueBlockChain(address string) *BlockChain {

	if !DBExists() {
		fmt.Println("No existing blockchain found, create a new one")
		runtime.Goexit()
	}

	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	utils.HandleError(err)
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.HandleError(err)
			err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		return err
	})
	utils.HandleError(err)
	bc := BlockChain{lastHash, db}
	return &bc
}

func (bc *BlockChain) FindUTXOutputs() map[string]TxOutputs {
	UTXOs := make(map[string]TxOutputs) 
	spentTxos := make(map[string][]int)
	iter := bc.Iterator()
	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTxos[txID] != nil {
					for _, spentOutIdx := range spentTxos[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}
				outs := UTXOs[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXOs[txID] = outs
		}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
						inID := hex.EncodeToString(in.ID)
						spentTxos[inID] = append(spentTxos[inID], in.OutIndex)
				}
			}
	}
		if len(block.PrevHash) == 0 {
			break
		}
	}	
	return UTXOs
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()
	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("transaction not found")

}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)
	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		utils.HandleError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privateKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {

	prevTXs := make(map[string]Transaction)
	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		utils.HandleError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}