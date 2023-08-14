package GRPC

import (
	BcGrpc "CoralDB/Proto/blockchain"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct{}

// 全局变量
//var LocalDataBlockChain *BCData.BlockChain
//var LocalTableBlockChain *BCTable.BlockChain
//var localTxPool *txpool.TxPool
//var LocalCache *cache.Cache
//
//func (s *Server) Init(LDbc *BCData.BlockChain, LTbc *BCTable.BlockChain, TxPool *txpool.TxPool, cache *cache.Cache) {
//	LocalDataBlockChain = LDbc
//	LocalTableBlockChain = LTbc
//	localTxPool = TxPool
//	LocalCache = cache
//	go s.startService()
//}

func (s *Server) Init() {
	//LocalDataBlockChain = LDbc
	//LocalTableBlockChain = LTbc
	//localTxPool = TxPool
	//LocalCache = cache
	go s.startService()
}

// 启动Grpc服务端
func (s *Server) startService() {
	// 初始化grpc对象
	grpcServer := grpc.NewServer()
	// 注册服务
	BcGrpc.RegisterBlockChainServiceServer(grpcServer, &Service{})
	// 创建监听
	listen, err := net.Listen("tcp", ":3301")
	if err != nil {
		log.Panic(err)
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Panic(err)
		}
	}(listen)

	// 绑定服务
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Panic(err)
	}
}
