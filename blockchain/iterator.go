package blockchain

import (
	"github.com/dgraph-io/badger"
	"github.com/nthskyradiated/blockchain-in-golang/utils"
)

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{bc.LastHash, bc.Database}
	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		utils.HandleError(err)
		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = append([]byte{}, val...)
			return nil
		})
		block = utils.Deserialize[*Block](encodedBlock)
		return err
	})
	utils.HandleError(err)
	iter.CurrentHash = block.PrevHash
	return block
}