package blockchain_data

import (
	"CoralDB/MerkleTree"
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func (blockChain *BlockChain) CheckDataBlock(block *Block) bool {
	//判断是否为创世区块
	if bytes.Equal(block.PreviousBlockHash, []byte("welcome to 407")) {
		fmt.Println("数据区块验证:  区块为创世区块")
		return true
	}
	// 校验交易是否合法
	for i := 0; i < len(block.Transactions); i++ {
		if !VerifyTransaction(*block.Transactions[i]) {
			fmt.Println("数据区块验证:  交易检验错误")
			return false
		}
	}
	// 校验默克尔根
	var MerKelRootData [][]byte
	// 类型转换
	for i := 0; i < len(block.Transactions); i++ {
		MerKelRootData = append(MerKelRootData, transactionsToBytes(*block.Transactions[i]))
	}
	MerKelRoot := MerkleTree.GetMerkleRoot(MerKelRootData).Hash
	if !bytes.Equal(MerKelRoot, block.MerKelRoot) {
		fmt.Println("数据区块验证:  默克尔根错误")
		log.Panic("数据区块验证:  默克尔根错误")
		return false
	}
	return true
}

// GetBlockByHash 通过hash得到区块
func (blockChain *BlockChain) GetBlockByHash(hash []byte) (*Block, error) {
	block := Block{}
	// 直接在数据库中查找
	err := blockChain.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockChain.BlockBucket))
		if bucket == nil {
			return errors.New("not found bucket")
		}
		value := bucket.Get(hash)
		if len(value) == 0 {
			return errors.New(fmt.Sprintf("not the key"))
		}
		block = Deserialize(value)
		return nil
	})
	if err != nil {
		if err.Error() == "not the key" {
			return nil, err
		} else {
			log.Panic(err)
		}
	}
	return &block, nil
}

// GetByRound 通过hash得到区块
func (blockChain *BlockChain) GetByRound(round uint64) *Block {
	lastBlock, err := blockChain.GetBlockByHash(blockChain.TailHash)
	if err != nil {
		log.Println("GetBlockById 获取区块失败")
	}
	if lastBlock.Round < round {
		return nil
	}
	var block *Block
	block = lastBlock
	for {
		if block.Round == round {
			return block
		}
		block1, err := blockChain.GetBlockByHash(block.PreviousBlockHash)
		if err != nil {
			return nil
		}
		block = block1
		if block.Round == 1 {
			break
		}
	}
	return nil
}
