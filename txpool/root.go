package txpool

import "fmt"

// 查看交易池里面的情况
func (tpl *TxPoolData) info() {
	fmt.Printf("队列长度: %d\n", tpl.txQueue.curSize)
	fmt.Printf("元素:")
	p := tpl.txQueue.head.next
	for p.next != nil {
		fmt.Printf("%s ", p.tx.Value)
		p = p.next
	}
}
