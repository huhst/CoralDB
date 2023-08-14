package cache

import (
	"fmt"
)

func (c *Cache) ROOT_PooledTables() {
	if len(c.tableInfo.tables) == 0 && len(c.tableInfo.tables) == 0 {
		fmt.Printf("没有任何表被加载\n")
		return
	}

	fmt.Println("(cache) 共享表用户权限:")
	for tableName, pl := range c.tableInfo.tables {
		fmt.Printf("table name : %s  ", tableName)
		for k, v := range pl {
			fmt.Printf("%s , %s;", k, v)
		}
		fmt.Println()
	}
}

func (c *Cache) ROOT_lruQ() {
	fmt.Printf("lru data len : %v; lru index len : %v\n", c.lru3Query.listTx.curSize, c.lru3Query.listIndex.curSize)
	fmt.Printf("map data len : %v; map index len : %v\n", c.lru3Query.listTx.curSize, c.lru3Query.listIndex.curSize)

	//p := c.lru3Query.listTx.head
	//ct1 := 0
	//for p != nil {
	//	ct1++
	//	p = p.next
	//}
	//fmt.Println("txLRU", ct1, " ", len(c.lru3Query.listTx.txMap))
	//fmt.Println()
	//
	//q := c.lru3Query.listIndex.head
	//ct2 := 0
	//for q != nil {
	//	ct2++
	//	q = q.next
	//}
	//fmt.Println("indexLRU", ct2, " ", len(c.lru3Query.listIndex.indexMap))
	//fmt.Println()
}
