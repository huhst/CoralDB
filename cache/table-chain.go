package cache

import (
	"CoralDB/blockchain/blockchain_data"
	"CoralDB/blockchain/blockchain_table"
	"CoralDB/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"sync"
)

// 记录表的信息。
// 1.共享表的权限 2.在tableDB里面建立共享表的相关链 3.提供查询接口

type tableInfo struct {
	sync.RWMutex
	tables     map[string]map[string]string // 表名，地址，权限
	dataChain  *blockchain_data.BlockChain
	tableChain *blockchain_table.BlockChain
	Possessor  map[string]string // 表名，创建者
	nodeNum    map[string]int    // 表名 ，协作者数量
}

// 初始化，读取所有的表。加载对应的权限信息到内存里面。
func (tio *tableInfo) init(dataChain *blockchain_data.BlockChain, tableChain *blockchain_table.BlockChain) {
	tio.Lock()
	defer tio.Unlock()

	tio.dataChain = dataChain
	tio.tableChain = tableChain

	tio.tables = make(map[string]map[string]string)
	tio.Possessor = make(map[string]string)
	tio.nodeNum = make(map[string]int)

	it := tio.tableChain.CreateIterator()
	for {
		// 放在前面可以不读创世区块。
		block := it.Next()
		if bytes.Equal(it.CurrentHash, []byte("welcome to 407")) == true {
			fmt.Printf("(cache ) : pooled tables Initialization complete\n")
			return
		}

		for _, tx := range block.Transactions {
			// 只读取最新的权限表信息
			if _, has := tio.tables[tx.Table]; !has {
				tio.tables[tx.Table] = make(map[string]string)
				tio.Possessor[tx.Table] = tx.Possessor
				nodeNum := 0
				for _, v := range tx.PermissionTable {
					tio.tables[tx.Table][v[:len(v)-1]] = v[len(v)-1:]
					nodeNum++
					tio.nodeNum[tx.Table] = nodeNum
				}
			}
		}
	}
}

// 创建新的表时，为新的表创建同名的bucket并初始化相关链。
func (tio *tableInfo) createBucket(tableName string) {

	tio.tableChain.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(tableName))
		if err != nil {
			log.Panic(err)
		}

		var key uint64
		keyBytes := util.Uint64ToBytes(key)
		err = bucket.Put(keyBytes, []byte("root"))

		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte(tableName), keyBytes)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
}

// 根据新的权限表的信息，更新内存中的权限表。
func (tio *tableInfo) upDateByTables(block blockchain_table.Block) {
	tio.Lock()
	defer tio.Unlock()

	for _, tx := range block.Transactions {

		if _, has := tio.tables[tx.Table]; !has {
			tio.tables[tx.Table] = make(map[string]string)
			tio.createBucket(tx.Table)
		}

		for _, v := range tx.PermissionTable {
			tio.tables[tx.Table][v[:len(v)-1]] = v[len(v)-1:]
		}
	}
}

// 在得到新的数据时，链长表相关链
func (tio *tableInfo) updateBucket(tableName string, blockHash []byte) {

	tio.tableChain.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		if bucket == nil {
			fmt.Println("(cache): 更新共享表相关链失败; 读取到空的bucket")
			os.Exit(1)
		}
		// 目标区块hash是否已近保存
		lastKeyBytes := bucket.Get([]byte(tableName))
		lastBlockHash := bucket.Get(lastKeyBytes)
		if bytes.Equal(lastBlockHash, blockHash) {
			return nil
		}

		// 添加新的区块hash到相关链
		lastKey := util.BytesToUint64(lastKeyBytes)
		lastKeyBytes = util.Uint64ToBytes(lastKey + 1)

		err := bucket.Put(lastKeyBytes, blockHash)
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte(tableName), lastKeyBytes)

		if err != nil {
			log.Panic(err)
		}
		return nil

	})
}

func (tio *tableInfo) upDateByData(block blockchain_data.Block) {
	for _, tx := range block.Transactions {
		tio.updateBucket(tx.Table, block.CurrentBlockHash)
	}
}

func (tio *tableInfo) getBlock(blockHash []byte) (block blockchain_data.Block) {

	tio.dataChain.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tio.dataChain.BlockBucket))
		if bucket == nil {
			log.Panic("bucket nil")
		}

		// 在boltDB里面得到区块
		bytesBlock := bucket.Get(blockHash)
		block = blockchain_data.Deserialize(bytesBlock)

		return nil

	})
	return block
}

func (tio *tableInfo) getInTableHashChain(dataID, tableName string) (Tx blockchain_data.Transaction, blockHash []byte, txID []byte, e error) {

	var lastTimeStamp int64
	var lastTX *blockchain_data.Transaction

	tio.tableChain.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		if bucket == nil {
			log.Panic("bucket nil")
		}

		lastKeyBytes := bucket.Get([]byte(tableName))
		lastBlockHash := bucket.Get(lastKeyBytes)

		for !bytes.Equal(lastBlockHash, []byte("root")) {

			block := tio.getBlock(lastBlockHash)
			for _, tx := range block.Transactions {
				if bytes.Equal(tx.DataID, []byte(dataID)) {
					if tx.TimeStamp > lastTimeStamp {
						lastTimeStamp = tx.TimeStamp
						lastTX = tx
					}
					blockHash = block.CurrentBlockHash
				}
			}
			if lastTimeStamp > 0 {
				return nil
			}
			lastKey := util.BytesToUint64(lastKeyBytes)
			lastKeyBytes = util.Uint64ToBytes(lastKey - 1)
			lastBlockHash = bucket.Get(lastKeyBytes)

		}

		return nil

	})
	if lastTimeStamp > 0 {
		return *lastTX, blockHash, lastTX.TxID, nil
	}
	return blockchain_data.Transaction{}, nil, nil, errors.New("null")
}

func (tio *tableInfo) getHistoryInTableHashChain(dataID string, tableName string) (txs []blockchain_data.Transaction, e error) {
	tio.RLock()
	defer tio.RUnlock()

	tio.tableChain.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		if bucket == nil {
			log.Panic("bucket nil")
		}

		lastKeyBytes := bucket.Get([]byte(tableName))
		lastBlockHash := bucket.Get(lastKeyBytes)

		for !bytes.Equal(lastBlockHash, []byte("root")) {
			block := tio.getBlock(lastBlockHash)
			for _, tx := range block.Transactions {
				if bytes.Equal(tx.DataID, []byte(dataID)) {
					txs = append(txs, *tx)
				}
			}
			lastKey := util.BytesToUint64(lastKeyBytes)
			lastKeyBytes = util.Uint64ToBytes(lastKey - 1)
			lastBlockHash = bucket.Get(lastKeyBytes)
		}

		return nil

	})
	if len(txs) == 0 {
		return nil, errors.New("null")
	}
	return txs, nil
}

// 得到同表里面的所有数据（根据key，拿到最新的数据）
// 根据表相关链，遍历区块。
// 得到新的数据时，添加。重复的key，比较时间戳在添加。
func (tio *tableInfo) getTableData(tableName string) (txs map[string]*blockchain_data.Transaction) {
	tio.RLock()
	defer tio.RUnlock()

	txs = make(map[string]*blockchain_data.Transaction)

	tio.tableChain.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		if bucket == nil {
			log.Panic("bucket nil")
		}

		lastKeyBytes := bucket.Get([]byte(tableName))
		lastBlockHash := bucket.Get(lastKeyBytes)

		for !bytes.Equal(lastBlockHash, []byte("root")) {
			block := tio.getBlock(lastBlockHash)
			for _, tx := range block.Transactions {
				if _, has := txs[tx.Key]; !has {
					txs[tx.Key] = tx
				} else {
					if tx.TimeStamp > txs[tx.Key].TimeStamp {
						txs[tx.Key] = tx
					}
				}
			}
			lastKey := util.BytesToUint64(lastKeyBytes)
			lastKeyBytes = util.Uint64ToBytes(lastKey - 1)
			lastBlockHash = bucket.Get(lastKeyBytes)
		}

		return nil

	})
	return txs
}

// CheckPermission 查找指定地址用户的权限
// 返回用户在指定表的权限 或者 错误信息
func (tio *tableInfo) checkPermission(address, table string) (string, int, error) {
	tio.RLock()
	defer tio.RUnlock()

	if _, has := tio.tables[table]; !has {
		return "-1", 0, errors.New("no such table")
	}

	if _, has := tio.tables[table][address]; !has {
		return "-1", 0, nil
	}
	return tio.tables[table][address], tio.nodeNum[table], nil
}

// uerTables 查看用户的所有表
func (tio *tableInfo) uerTables(address string, isWeb bool) (myTables map[string]string, err error) {
	tio.RLock()
	defer tio.RUnlock()

	myTables = make(map[string]string)
	for tn, tpl := range tio.tables {
		if pl, has := tpl[address]; has {
			if isWeb {
				myTables[tn] = pl + "," + tio.Possessor[tn]
			} else {
				myTables[tn] = pl
			}
		}
	}

	if len(myTables) == 0 {
		return nil, errors.New("null")
	}
	return myTables, nil
}
