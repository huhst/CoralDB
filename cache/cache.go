package cache

import (
	"CoralDB/blockchain/blockchain_data"
	"CoralDB/blockchain/blockchain_table"
	"bytes"
	"errors"
	"fmt"
	"sync"
)

var LocalCache *Cache

// cache
// 1. 表的权限信息  2. LRU-3 的缓存队列  3. 表相关链条
// 这些缓存信息。在程序启动时要读取区块链文件初始化，在本地接受到新的区块时，要读取新的区块来更新缓存。
// 读取数据区块更新 lru3Query。 读取表区块更新 powerTable。 数据区块或者表区块都要更新的是 表相关链 tableHashChain。
// funcAPI 1. 在程序启动时，初始化
// funcAPI 2. 在拿到新的区块时，更新
// funcAPI 3. 对外的查询接口 （1.在 lru 中查找， 2.在表相关链中查找  3.遍历boltDB查找（原则是是不需要的）
// funcAPI 其他。
// a. CheckPermission(address, table string) 查看指定地址在表里面的权限

type Cache struct {
	sync.Mutex

	lru3Query lru3Query
	tableInfo tableInfo

	dataChain  *blockchain_data.BlockChain
	tableChain *blockchain_table.BlockChain
}

func (c *Cache) Init(dc *blockchain_data.BlockChain, tc *blockchain_table.BlockChain) {
	c.Lock()
	defer c.Unlock()

	c.dataChain = dc
	c.tableChain = tc
	c.lru3Query.init(dc)
	c.tableInfo.init(dc, tc)
	LocalCache = c
}

func (c *Cache) GetOneValue(dataID string, tableName string) (blockchain_data.Transaction, error) {
	// 1.在lru3Query 数据队列 里面查找元素
	tx, err := c.lru3Query.getInTX(dataID)
	if err == nil {
		c.lru3Query.listTx.toTail(dataID)
		//fmt.Printf("find in cache  ")
		return tx, nil
	}
	// 2.在lru3Query 索引队列 里面查找元素
	tx, err = c.lru3Query.getInIndex(dataID)
	if err == nil { // 找到了具体的数据
		c.lru3Query.listIndex.toTail(dataID)
		if c.lru3Query.listIndex.checkDel(dataID) {
			c.lru3Query.listTx.queueIn(dataID, tx)
		}
		//dc.CacheInfo()
		//fmt.Printf("find in index  ")
		return tx, nil
	}
	// 3. 在表的相关链里面查找
	tx, blockHash, txID, err := c.tableInfo.getInTableHashChain(dataID, tableName)
	if err == nil {
		c.lru3Query.listIndex.queueIn(dataID, blockHash, txID)
		//dc.CacheInfo()
		//fmt.Printf("find in table hash chain  ")
		return tx, nil
	}

	// 4. 在boltDB 里面查找
	tx, blockHash, txID, err = c.getInBoltDb(dataID)
	if err == nil {
		c.lru3Query.listIndex.queueIn(dataID, blockHash, txID)
		//dc.CacheInfo()
		//fmt.Printf("find in boltDB  ")
		return tx, nil
	}

	return tx, errors.New("null")
}

// 在boltDB里面查找数据, 返回交易或者 err.
// 查找的最后一种可能，遍历区块链查找数据。
// 处理一种情况。一个区块里面有对一个数据的多个交易时那么就需要得到这个最新时间的数据
func (c *Cache) getInBoltDb(ID string) (blockchain_data.Transaction, []byte, []byte, error) {
	it := c.dataChain.CreateIterator()
	// 通过时间戳确定同个区块，同ID的最新数据
	var lastTimeStamp int64
	var lastTX *blockchain_data.Transaction

	for {
		// 遍历区块
		block := it.Next()
		if bytes.Equal(it.CurrentHash, []byte("welcome to 407")) == true {
			return blockchain_data.Transaction{}, nil, nil, errors.New("null")
		}

		// 遍历交易
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.DataID, []byte(ID)) {
				fmt.Println("here")
				if tx.TimeStamp > lastTimeStamp {
					lastTimeStamp = tx.TimeStamp
					lastTX = tx
				}
			}
		}
		// 在当前区块里面找到了数据，直接返回。
		if lastTimeStamp > 0 {
			return *lastTX, block.CurrentBlockHash, lastTX.TxID, nil
		}
	}
}

// GetHistoryValue 查找指定数据的历史记录
func (c *Cache) GetHistoryValue(dataID string, tableName string) (txs []blockchain_data.Transaction, e error) {
	return c.tableInfo.getHistoryInTableHashChain(dataID, tableName)
}

func (c *Cache) GetTableData(tableName string) map[string]*blockchain_data.Transaction {
	return c.tableInfo.getTableData(tableName)
}

// UpdateByDataBlock 在接受新的数据区块时，更新缓存中的 lru3Query 和 tableHashChain
func (c *Cache) UpdateByDataBlock(block blockchain_data.Block) {
	c.lru3Query.updateDataCache(block)
	c.tableInfo.upDateByData(block)
}

// UpdateByTableBlock 在接受新的共享表区块时， 更新缓存中的 powerTable 和 tableHashChain
func (c *Cache) UpdateByTableBlock(block blockchain_table.Block) {
	c.tableInfo.upDateByTables(block)
}

// CheckPermission 返回对应地址在指定表的权限
func (c *Cache) CheckPermission(address, tableName string) (string, int, error) {
	return c.tableInfo.checkPermission(address, tableName)
}

func (c *Cache) MyTables(address string, isWeb bool) (myTables map[string]string, err error) {
	return c.tableInfo.uerTables(address, isWeb)
}
