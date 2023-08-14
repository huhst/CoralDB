package txpool

import (
	"CoralDB/blockchain/blockchain_data"
	"CoralDB/blockchain/blockchain_table"
	"CoralDB/cache"
	"CoralDB/util"
	"fmt"
	"sync"
	"time"
)

var LocalTxPool *TxPool

// TxPool TxPool: 交易池，根据本地节点的状态来对交易池里面的交易处理。
// 交易池定期检查自己在集群中的状态，如果是记账节点，执行记账行动。
// 交易池接受被校验过的交易进入交易池。
// tip 同时维护了数据链和数据共享链的交易池，在输入输出和对应更新交易池时要注意区分。
// 互斥锁。 1.在定期执行记账动作时 2.修改交易池的设置  tip：对交易池队列的锁，下放到具体的交易池队列里面。
// funcAPI
// 1. setMod 修改交易池的处理模式，如果是修改为记账节点。开始记账活动。
// 2. SetPackNumber 修改交易池打包区块时，区块里面交易的数量。
// 3. TxDataIN， TxTableN 交易进入对应的交易池。
// 4. UpdateByData UpdateByTable 普通节点得到记账节点的区块时，用里面的交易来更新自己的交易池。

type TxPool struct {
	sync.RWMutex

	tick       *time.Ticker
	bookKeeper bool

	txPoolData  TxPoolData
	txPoolTable TxPoolTable
}

func (tpl *TxPool) Init(chain1 *blockchain_data.BlockChain, chain2 *blockchain_table.BlockChain, cache *cache.Cache) {
	tpl.tick = time.NewTicker(time.Millisecond * 200)
	tpl.bookKeeper = false
	tpl.txPoolData.init(chain1, cache)
	tpl.txPoolTable.init(chain2, cache)
	LocalTxPool = tpl
	go tpl.run()

	fmt.Printf("(tx pool ) : Initialization complete\n")
}

func (tpl *TxPool) run() {
	defer tpl.tick.Stop()
	// 定期检查本地节点的状态
	for {
		//if util.IsDone {
		//fmt.Println()
		tpl.txPoolData.bookKeeperRun()
		//}
		<-tpl.tick.C
		tpl.Lock()
		//if tpl.bookKeeper{
		if util.LocalIP == "192.169.10.2" {
			tpl.txPoolTable.bookKeeperRun()
			//tpl.txPoolData.bookKeeperRun()
		}
		tpl.Unlock()
	}
}

func (tpl *TxPool) SetMod(isAccount bool) {
	tpl.Lock()
	defer tpl.Unlock()
	if tpl.bookKeeper == isAccount {
		return
	}
	tpl.bookKeeper = !tpl.bookKeeper
	//if tpl.bookKeeper == true {
	//	fmt.Printf("设置当前节点为记账节点; 数据区块打包的交易数量为 %d ;\n", tpl.txPoolData.packNumber)
	//} else {
	//	fmt.Printf("设置当前节点为普通节点;\n")
	//}
}
func (tpl *TxPool) SetMod1() {
	tpl.Lock()
	defer tpl.Unlock()

	tpl.bookKeeper = !tpl.bookKeeper
	if tpl.bookKeeper == true {
		fmt.Printf("设置当前节点为记账节点; 数据区块打包的交易数量为 %d ;\n", tpl.txPoolData.packNumber)
	} else {
		fmt.Printf("设置当前节点为普通节点;\n")
	}
}

func (tpl *TxPool) SetPackNumber(num int) {
	tpl.Lock()
	defer tpl.Unlock()

	tpl.txPoolData.packNumber = num
	// 共享表区块链不作修改，默认就是一个表，一个区块。
	// tpl.txPoolTable.packNumber = 1
	fmt.Printf("修改区块打包交易数量为 %d;\n", tpl.txPoolData.packNumber)
}

func (tpl *TxPool) TxDataIN(transaction blockchain_data.Transaction) error {
	err := tpl.txPoolData.TxIn(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (tpl *TxPool) TxTableIN(transaction blockchain_table.Transaction) error {
	err := tpl.txPoolTable.TxIn(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (tpl *TxPool) UpdateByData(block blockchain_data.Block) {
	tpl.txPoolData.OrdinaryRun(block)
}

func (tpl *TxPool) UpdateByTable(block blockchain_table.Block) {
	tpl.txPoolTable.OrdinaryRun(block)
}
