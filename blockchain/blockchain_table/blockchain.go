package blockchain_table

import (
	"CoralDB/util"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

type BlockChain struct {
	Db       *bolt.DB // Db 本地的boltDB 连接
	TailHash []byte   // TailHash 区块链最后一个区块的 HASH

	BlockChainFile string
	BlockBucket    string
	LastHashKey    string
	LastID         int
}

var LocalTableBlockChain *BlockChain

// InitBlockChain 如果本地存在区块区块链文件，就更新数据
// 如果本地不存在区块区块链文件，创建新的区块链。加入genesisBlock.
func (blockChain *BlockChain) InitBlockChain() {
	blockChain.BlockChainFile = "tableBlockChain.db"
	blockChain.BlockBucket = "blockBucket"
	blockChain.LastHashKey = "lastHashKey"
	blockChain.LastID = 1

	if !util.IsExistFile(blockChain.BlockChainFile) {
		fmt.Println("本地不存在表权限区块链文件, 创建区块链文件 ing .")
		// 创建 boltDB
		db, err := bolt.Open(blockChain.BlockChainFile, 0600, nil)
		if err != nil {
			log.Panic(err)
		}
		// 创建 boltDB.bucket
		db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucket([]byte(blockChain.BlockBucket))
			if err != nil {
				log.Panic(err)
			}
			// 创建genesisBlock
			// "welcome to 407" 创世区块的HASH
			genesisBlock := NewBlock()
			transaction := Transaction{Table: "sample table"}
			genesisBlock.InitGenesisBlock([]*Transaction{&transaction}, []byte("welcome to 407"), blockChain.LastID)
			blockChain.LastID++
			// 更新区块链
			err = bucket.Put(genesisBlock.CurrentBlockHash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put([]byte(blockChain.LastHashKey), genesisBlock.CurrentBlockHash)
			if err != nil {
				log.Panic(err)
			}

			// Init blockChain
			blockChain.Db = db
			blockChain.TailHash = genesisBlock.CurrentBlockHash
			return nil
		})
	} else {
		fmt.Println("加载区块链文件 ing .")
		// 链接 boltDB
		db, err := bolt.Open(blockChain.BlockChainFile, 0600, nil)
		if err != nil {
			log.Panic(err)
		}
		// 链接 boltDB.bucket
		db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(blockChain.BlockBucket))
			if bucket == nil {
				fmt.Println("BlockBucket is nil")
				os.Exit(1)
			}

			// Init blockchain
			blockChain.Db = db
			blockChain.TailHash = bucket.Get([]byte(blockChain.LastHashKey))
			return nil
		})
	}
	LocalTableBlockChain = blockChain
	fmt.Println("区块链初始化完成.")
}

// AddBlockToChain 添加区块到区块链
func (blockChain *BlockChain) AddBlockToChain(block Block) {
	blockChain.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockChain.BlockBucket))
		if bucket == nil {
			fmt.Println("BlockBucket is nil")
			os.Exit(1)
		}
		// 添加区块并更新信息
		err := bucket.Put(block.CurrentBlockHash, block.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte(blockChain.LastHashKey), block.CurrentBlockHash)
		if err != nil {
			log.Panic(err)
		}

		blockChain.TailHash = block.CurrentBlockHash
		return nil
	})
}

// Iterator 区块链迭代器
type Iterator struct {
	db          *bolt.DB
	BlockBucket string
	CurrentHash []byte
}

func (blockChain *BlockChain) CreateIterator() Iterator {
	return Iterator{db: blockChain.Db, CurrentHash: blockChain.TailHash, BlockBucket: blockChain.BlockBucket}
}

func (iterator *Iterator) Next() Block {
	var block Block
	iterator.db.View(func(tx *bolt.Tx) error {
		//fmt.Println("迭代器hash", string(iterator.currentHash))

		bucket := tx.Bucket([]byte(iterator.BlockBucket))
		if bucket == nil {
			fmt.Println("BlockBucket is nil")
			os.Exit(1)
		}
		// 在boltDB里面得到区块
		bytesBlock := bucket.Get(iterator.CurrentHash)

		block = Deserialize(bytesBlock)
		iterator.CurrentHash = block.PreviousBlockHash

		return nil
	})
	return block
}
