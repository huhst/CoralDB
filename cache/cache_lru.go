package cache

import (
	"CoralDB/blockchain/blockchain_data"
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"sync"
)

type txNode struct {
	tx        blockchain_data.Transaction
	pre, next *txNode
}

type indexNode struct {
	blockHash []byte
	txID      []byte
	count     int
	dataID    string

	pre, next *indexNode
}

type lruTx struct {
	maxSize int
	curSize int

	txMap map[string]*txNode
	head  *txNode
	tail  *txNode
}

type lruIndex struct {
	maxSize int
	curSize int

	indexMap map[string]*indexNode
	head     *indexNode
	tail     *indexNode
}

func (lt *lruTx) init() {
	lt.maxSize = 120
	lt.curSize = 0
	lt.txMap = make(map[string]*txNode)
	p1 := new(txNode)
	p2 := new(txNode)
	lt.head = p1
	lt.tail = p2
	p1.next = p2
	p2.pre = p1
}

func (li *lruIndex) init() {
	li.maxSize = 300
	li.curSize = 0
	li.indexMap = make(map[string]*indexNode)
	p1 := new(indexNode)
	p2 := new(indexNode)
	li.head = p1
	li.tail = p2
	p1.next = p2
	p2.pre = p1
}

type lru3Query struct {
	sync.RWMutex

	listTx    lruTx
	listIndex lruIndex
	chain     *blockchain_data.BlockChain
}

func (lru *lru3Query) init(chain *blockchain_data.BlockChain) {
	lru.Lock()
	defer lru.Unlock()

	lru.listTx.init()
	lru.listIndex.init()
	lru.chain = chain
}

// 维护队列之在队列里面插入元素
// 维护缓存交易队列
func (lt *lruTx) queueIn(ID string, tx blockchain_data.Transaction) {
	p := &txNode{
		tx:   tx,
		pre:  nil,
		next: nil,
	}

	if lt.curSize > lt.maxSize {
		pg := lt.head.next
		pg.next.pre = pg.pre
		pg.pre.next = pg.next
		delete(lt.txMap, string(pg.tx.DataID))
		lt.curSize--
	}

	ptag := lt.tail.pre
	p.next = ptag.next
	ptag.next.pre = p
	p.pre = ptag
	ptag.next = p

	lt.txMap[ID] = p
	lt.curSize++
	return
}

// 维护缓存索引队列
func (li *lruIndex) queueIn(ID string, blockHash, txID []byte) {
	p := &indexNode{
		blockHash: blockHash,
		txID:      txID,
		dataID:    ID,
		count:     1, // 加入队列即为被最近被访问一次
		pre:       nil,
		next:      nil,
	}
	if li.curSize > li.maxSize {
		pg := li.head.next
		pg.next.pre = pg.pre
		pg.pre.next = pg.next
		delete(li.indexMap, pg.dataID)
		li.curSize--
	}

	ptag := li.tail.pre
	p.next = ptag.next
	ptag.next.pre = p
	p.pre = ptag
	ptag.next = p

	li.indexMap[ID] = p
	li.curSize++

	return
}

// 更新指定的元素到交易队列尾部
func (lt *lruTx) toTail(ID string) {

	p := lt.txMap[ID]
	if lt.tail.pre == p {
		return
	}

	// 断开
	p.pre.next = p.next
	p.next = p.pre
	// 接上
	ptag := lt.tail.pre
	p.next = ptag.next
	ptag.next.pre = p
	p.pre = ptag
	ptag.next = p

}

// 更新指定的元素到索引队列尾部
func (li *lruIndex) toTail(ID string) {
	// 移动一次意味着访问次数加一
	p := li.indexMap[ID]
	p.count++

	// 在队尾，不用处理
	if li.tail.pre == p {
		return
	}

	// 断开
	p.pre.next = p.next
	p.next = p.pre
	// 接上
	ptag := li.tail.pre
	p.next = ptag.next
	ptag.next.pre = p
	p.pre = ptag
	ptag.next = p

	return
}

// checkDel 判断这个节点时候满足次数，如果满足就删除
// 应为已经是移动到了 tail，所以直接判断tail就可以了。
func (li *lruIndex) checkDel(ID string) bool {
	p := li.indexMap[ID]

	if p.count >= 3 {
		p.pre.next = p.next
		p.next.pre = p.pre

		delete(li.indexMap, ID)
		li.curSize--

		return true
	}
	return false
}

// 在交易队列里面查找，返回交易或者err
// 通过数据ID->这个数据最新的交易
func (lru *lru3Query) getInTX(dataID string) (blockchain_data.Transaction, error) {
	lru.Lock()
	defer lru.Unlock()

	if p, has := lru.listTx.txMap[dataID]; has {
		return p.tx, nil
	}
	return blockchain_data.Transaction{}, errors.New("null")
}

// 在索引队列里面查找，索引或者err
// 由数据ID-> blockHASH, txID 。根据得到的索引来直接在boltDB里面找交易，而不是遍历区块链来找到交易
// 可以直接在boltDB里面定位到具体的交易
func (lru *lru3Query) getInIndex(dataID string) (blockchain_data.Transaction, error) {
	lru.Lock()
	defer lru.Unlock()

	var tx0 blockchain_data.Transaction

	if p, has := lru.listIndex.indexMap[dataID]; has {
		blockHash := p.blockHash
		txID := p.txID

		var block blockchain_data.Block
		lru.chain.Db.View(func(tx *bolt.Tx) error {

			bucket := tx.Bucket([]byte(lru.chain.BlockBucket))
			if bucket == nil {
				fmt.Println("BlockBucket is nil")
				os.Exit(1)
			}
			// 在boltDB里面得到区块
			bytesBlock := bucket.Get(blockHash)
			block = blockchain_data.Deserialize(bytesBlock)

			for _, tx := range block.Transactions {
				if bytes.Equal(tx.TxID, txID) {
					tx0 = *tx
				}
			}
			return nil
		})
		return tx0, nil
	}
	return tx0, errors.New("null")
}

// UpdateDataCache 得到新区块的时候 更新缓存
func (lru *lru3Query) updateDataCache(block blockchain_data.Block) {
	lru.Lock()
	defer lru.Unlock()

	// 因为是遍历了所有的交易，那么得到的信息就是区块里面靠后的。
	for _, tx := range block.Transactions {
		dataId := string(tx.DataID)
		if _, has := lru.listTx.txMap[dataId]; has {
			lru.listTx.txMap[dataId].tx = *tx
		}
		if _, has := lru.listIndex.indexMap[dataId]; has {
			lru.listIndex.indexMap[dataId].blockHash = block.CurrentBlockHash
			lru.listIndex.indexMap[dataId].txID = tx.TxID
		}
	}
}
