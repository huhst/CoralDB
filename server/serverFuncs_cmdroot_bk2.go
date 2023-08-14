package server

import (
	"fmt"
	"time"
)

func (s *Server) GetTableData(UID, table string) {
	start := time.Now() // 获取当前时间

	// 登录状态检查
	a, err := s.manage.ViewAccount(UID)
	if err != nil {
		fmt.Println("用户未登录")
		return
	}
	ad := a.Address

	// 权限检查
	permission,_, err := s.Cache.CheckPermission(ad, table)
	if err != nil {
		fmt.Printf("不存在这个表; %s\n", table)
		return
	} else {
		if permission < "1" {
			fmt.Println("没有对表的查看权限")
			return
		}
	}

	txs := s.Cache.GetTableData(table)
	if len(txs) == 0 {
		fmt.Printf("表为空,里面没有任何数据")
	}
	for _, tx := range txs {
		fmt.Printf("key : %s    ", tx.Key)
		fmt.Printf("value: %s    ", tx.Value)
		fmt.Printf("possessor: %s    ", tx.Possessor)
		fmt.Printf("alterTime : %v\n", time.Unix(tx.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
	}
	elapsed := time.Since(start)
	fmt.Println("该查找执行完成耗时：", elapsed)
}
