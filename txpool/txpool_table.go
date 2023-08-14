package txpool

import (
	"CoralDB/blockchain/blockchain_table"
	"CoralDB/blockqueue"
	"CoralDB/cache"
	"errors"
	"fmt"
	"sync"
	"time"
)

type tableTxNode struct {
	tx   blockchain_table.Transaction
	pre  *tableTxNode
	next *tableTxNode
}

type tableQueue struct {
	curSize int
	maxSize int
	txMap   map[string]*tableTxNode // map[tableName]*tableTxNode
	head    *tableTxNode
	tail    *tableTxNode
}

func (q *tableQueue) init() {
	q.curSize = 0
	q.maxSize = 10

	q.txMap = make(map[string]*tableTxNode)
	p1 := new(tableTxNode)
	p2 := new(tableTxNode)
	q.head = p1
	q.tail = p2
	p1.next = p2
	p2.pre = p1
}

func (q *tableQueue) in(transaction blockchain_table.Transaction) {
	p := &tableTxNode{
		tx:   transaction,
		pre:  nil,
		next: nil,
	}

	// 找到合适的位置, 自己的时间戳是要大于前面节点的时间戳的
	tsp := p.tx.TimeStamp
	p1 := q.tail.pre
	for p1.pre != nil && p1.tx.TimeStamp > tsp { // p1 不是头节点，并且p1的时间靠后
		p1 = p1.pre
	} // result p1 时间靠前，或者是头节点

	p.next = p1.next
	p1.next.pre = p
	p.pre = p1
	p1.next = p

	q.txMap[string(p.tx.TxID)] = p
	q.curSize++
}

func (q *tableQueue) out(transaction blockchain_table.Transaction) {

	p := q.txMap[string(transaction.TxID)]
	p.pre.next = p.next
	p.next.pre = p.pre

	delete(q.txMap, string(transaction.TxID))
	q.curSize--
}

// 记账节点的动作。 定时在等待队列中拿取打包

type TxPoolTable struct {
	sync.Mutex

	tableQueue tableQueue

	packNumber   int
	count        int
	countPointer *tableTxNode

	chain *blockchain_table.BlockChain
	cache *cache.Cache
}

func (tpl *TxPoolTable) init(chain *blockchain_table.BlockChain, cache *cache.Cache) {
	tpl.Lock()
	defer tpl.Unlock()

	tpl.tableQueue.init()

	tpl.packNumber = 1
	tpl.count = 0
	tpl.countPointer = tpl.tableQueue.head
	tpl.chain = chain
	tpl.cache = cache
}

func (tpl *TxPoolTable) bookKeeperRun() {
	tpl.Lock()
	defer tpl.Unlock()

	tspNow := time.Now().Unix()
	tspStand := tspNow - 2 //时间戳2s前

	// 检查等待队列的交易的时间戳，统计时间合法的交易。
	// 从head开始计数时间戳符合要求的节点数量
	tpl.count = 0
	tpl.countPointer = tpl.tableQueue.head.next
	for tpl.countPointer.next != nil && tpl.countPointer.tx.TimeStamp < tspStand {
		tpl.countPointer = tpl.countPointer.next
		tpl.count++
	} // result: cp-> 尾节点，或者第一个时间戳没有达到延迟要求的节点

	// 如果数量满足打包的要求，将这些交易批量打包。
	for tpl.count >= tpl.packNumber {

		var txs []*blockchain_table.Transaction
		p := tpl.tableQueue.head.next // start: p.第一个要被打包的元素
		for i := 0; i < tpl.packNumber; i++ {
			txs = append(txs, &p.tx)
			delete(tpl.tableQueue.txMap, string(p.tx.TxID))
			tpl.tableQueue.curSize--

			p = p.next
		} // result: p指向一个不被打包的元素。或者刚好打包完，指向尾节点。

		block := blockchain_table.NewBlock()
		block.InitBlock(txs, tpl.chain.TailHash, tpl.chain.LastID)
		tpl.chain.LastID++
		tpl.chain.AddBlockToChain(block)
		fmt.Println("生成一个新的共享表区块", time.Now().String())

		// 更新本地缓存
		tpl.cache.UpdateByTableBlock(block)

		// TODO 分发区块
		//fmt.Println("区块入队")
		blockqueue.LocalTableBlockQueue.Put(block)
		//fmt.Println("入队完成")
		tpl.tableQueue.head.next = p
		p.pre = tpl.tableQueue.head

		tpl.count -= tpl.packNumber
	}
}

// OrdinaryRun 普通节点在接受到记账节点，发过来的区块时。根据区块里面的交易，删除自己交易池里面对应的交易。
func (tpl *TxPoolTable) OrdinaryRun(block blockchain_table.Block) {
	tpl.Lock()
	defer tpl.Unlock()

	for _, tx := range block.Transactions {
		tpl.tableQueue.out(*tx)
	}
}

// TxIn 被校验过的交易进入交易池
// 对时间戳进行检查，找到合适的位置并插入到队列里面。
// 记账节点： 接受其他节点和自己的交易进入交易池
// 普通节点： 接受记账节点传过来的交易然后进入交易池
func (tpl *TxPoolTable) TxIn(transaction blockchain_table.Transaction) error {
	tpl.Lock()
	defer tpl.Unlock()

	if tpl.tableQueue.curSize == tpl.tableQueue.maxSize {
		return errors.New("FULL")
	}

	tpl.tableQueue.in(transaction)

	return nil
}
