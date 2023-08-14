package server

import (
	"CoralDB/GRPC"
	"CoralDB/blockchain/blockchain_data"
	"CoralDB/blockchain/blockchain_table"
	"CoralDB/cache"
	"CoralDB/txpool"
	"CoralDB/userManage"
	"sync"
)

// Server 定义一个区块链数据库服务
type Server struct {
	mutex      sync.Mutex
	manage     userManage.UserManage
	dataChain  blockchain_data.BlockChain
	tableChain blockchain_table.BlockChain
	TxPool     txpool.TxPool
	Cache      cache.Cache
	Grpc       GRPC.Server
}

func (s *Server) Init() {
	s.manage.Init()
	s.dataChain.InitBlockChain()
	s.tableChain.InitBlockChain()
	s.Cache.Init(&s.dataChain, &s.tableChain)
	s.TxPool.Init(&s.dataChain, &s.tableChain, &s.Cache)
	s.Grpc.Init()
}
