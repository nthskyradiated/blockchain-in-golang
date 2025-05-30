package blockchain

import (
	"fmt"
	"os"	
	"github.com/dgraph-io/badger"
	"github.com/nthskyradiated/blockchain-in-golang/utils"
)

const (
	dbPath = "./tmp/blocks"
	genesisData = "Genesis Block"
)
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func DBExists(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}
	return true
}

func (bc *BlockChain) AddBlock(data string) {

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


	newBlock := CreateBlock(data, lastHash)

	err = bc.Database.Update(func(txn *badger.Txn) error {
	err := txn.Set(newBlock.Hash, newBlock.Serialize())
	utils.HandleError(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return err
	})
	utils.HandleError(err)
}

func NewBlockChain() *BlockChain {
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	utils.HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No last hash found, creating genesis block")
			genesis := GenesisBlock()
			fmt.Println("Genesis block created")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			utils.HandleError(err)
			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash
			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			utils.HandleError(err)
			err = item.Value(func (val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
			return err
		}
	})
	utils.HandleError(err)
	bc := BlockChain{lastHash, db}
	return &bc
}