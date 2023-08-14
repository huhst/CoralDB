package GRPC

import (
	"CoralDB/Cluster"
	BcGrpc "CoralDB/Proto/blockchain"
	"CoralDB/consensus"
	BCData "CoralDB/blockchain/blockchain_data"
	BCTable "CoralDB/blockchain/blockchain_table"
	"CoralDB/cache"
	"CoralDB/txpool"
	"CoralDB/util"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
)

// Service GRPC接口的实现(客户端)
type Service struct{}

// DistributeDataBlock 数据区块的分发
func (s *Service) DistributeDataBlock(ctx context.Context, req *BcGrpc.DataBlock) (info *BcGrpc.VerifyInfo, err error) {
	util.IsDone = true
	info = &BcGrpc.VerifyInfo{}
	info.Info = "区块校验成功"
	info.Status = true
	// 区块的转换
	block := GrpcDataBlockToBlock(req)
	// 区块的校验
	//fmt.Println("rpcService 数据区块的进入校验", block)
	flag := BCData.LocalDataBlockChain.CheckDataBlock(block)
	if !flag {
		info.Info = "区块校验失败"
		info.Status = false
	}
	// 将区块上链
	BCData.LocalDataBlockChain.AddBlockToChain(*block)
	// 交易池更新
	// 不是提议节点才进行更新
	if util.SubUser <= 0 {
		//log.Println("交易池更新")
		txpool.LocalTxPool.UpdateByData(*block)
	}
	// 缓存更新
	//log.Println("缓存更新")
	cache.LocalCache.UpdateByDataBlock(*block)
	//log.Println("完成")
	return info, nil
}

// DataBlockSynchronization 数据区块的同步
func (s *Service) DataBlockSynchronization(ctx context.Context, req *BcGrpc.ReqDataBlock) (out *BcGrpc.ResDataBlocks, err error) {
	out = &BcGrpc.ResDataBlocks{}

	// 获得迭代器
	iterator := BCData.LocalDataBlockChain.CreateIterator()
	// 获得未同步的区块
	for {
		block := iterator.Next()
		if !bytes.Equal(block.CurrentBlockHash, req.Hash) {
			out.Blocks = append(out.Blocks, BlockToGrpcDataBlock(&block))
		} else {
			break
		}
	}
	return out, nil
}

// DistributeTableBlock 表区块的分发
func (s *Service) DistributeTableBlock(ctx context.Context, req *BcGrpc.TableBlock) (info *BcGrpc.VerifyInfo, err error) {
	info = &BcGrpc.VerifyInfo{}
	info.Info = "区块校验成功"
	info.Status = true
	// 区块的转换
	block := GrpcTableBlockToBlock(req)
	// 区块的校验
	//fmt.Println("rpcService 权限区块的进入校验")
	flag := BCTable.LocalTableBlockChain.CheckTableBlock(block)
	if !flag {
		info.Info = "区块校验失败"
		info.Status = false
	}
	// 上链
	BCTable.LocalTableBlockChain.AddBlockToChain(*block)
	// 交易池更新
	txpool.LocalTxPool.UpdateByTable(*block)
	// 缓存更新
	cache.LocalCache.UpdateByTableBlock(*block)
	//log.Println("表区块的分发")
	return info, nil
}

// TableBlockSynchronization 表区块的同步
func (s *Service) TableBlockSynchronization(ctx context.Context, req *BcGrpc.ReqTableBlock) (out *BcGrpc.ResTableBlocks, err error) {
	out = &BcGrpc.ResTableBlocks{}
	// 获得Leader上的区块链

	// 获得迭代器
	iterator := BCTable.LocalTableBlockChain.CreateIterator()
	// 获得未同步的区块
	for {
		block := iterator.Next()
		if !bytes.Equal(block.CurrentBlockHash, req.Hash) {
			out.Blocks = append(out.Blocks, BlockToGrpcTableBlock(&block))
		} else {
			break
		}
	}
	return out, nil
}

//// NodeHeartbeat 心跳检测
//func (s *Server) NodeHeartbeat(ctx context.Context, req *BcGrpc.NodeInfo) (out *BcGrpc.Heartbeat, err error) {
//
//	return nil, status.Errorf(codes.Unimplemented, "method NodeHeartbeat not implemented")
//}
//
//// GetAccountant 获取记账权
//func (s *Server) GetAccountant(ctx context.Context, req *BcGrpc.Heartbeat) (info *BcGrpc.VerifyInfo, err error) {
//
//	return nil, status.Errorf(codes.Unimplemented, "method GetAccountant not implemented")
//}

// DataTradingPool  数据交易池
func (s *Service) DataTradingPool(ctx context.Context, req *BcGrpc.DataTransaction) (info *BcGrpc.VerifyInfo, err error) {
	info = &BcGrpc.VerifyInfo{}
	info.Info = "已接收交易"
	info.Status = true
	// 交易的转换
	tx := GrpcDataTxToDataTx(req)
	// 验证交易
	if !BCData.VerifyTransaction(*tx) {
		info.Info = "交易验证失败"
		info.Status = false
		return info, errors.New("交易验证失败")
	}
	// 交易入本地交易池
	err = txpool.LocalTxPool.TxDataIN(*tx)
	if err != nil {
		info.Info = "交易入池失败"
		info.Status = false
		return info, err
	}
	return info, nil
}

// TableTradingPool  表交易池
func (s *Service) TableTradingPool(ctx context.Context, req *BcGrpc.TableTransaction) (info *BcGrpc.VerifyInfo, err error) {
	//fmt.Println("调用交易入池")
	info = &BcGrpc.VerifyInfo{}
	info.Info = "已接收交易"
	info.Status = true
	// 交易的转换
	tx := GrpcTableTxToTableTx(req)
	// 验证交易
	if !BCTable.VerifyTransaction(*tx) {
		info.Info = "交易验证失败"
		info.Status = false
		return info, errors.New("交易验证失败")
	}
	// 交易入本地交易池
	err = txpool.LocalTxPool.TxTableIN(*tx)
	if err != nil {
		info.Info = "交易入池失败"
		info.Status = false
		return info, err
	}
	return info, nil
}

// DataTradingPoolTxs 多交易的广播
func (s *Service) DataTradingPoolTxs(ctx context.Context, req *BcGrpc.DataTransactions) (info *BcGrpc.VerifyInfo, err error) {
	//start := time.Now().UnixMilli()
	//fmt.Println("ip ", util.LocalIP, "req in ", time.UnixMicro(time.Now().UnixMicro()))
	info = &BcGrpc.VerifyInfo{}
	info.Info = "已接收交易"
	info.Status = true
	for i := 0; i < len(req.Transactions); i++ {
		tx := GrpcDataTxToDataTx(req.Transactions[i])
		if !BCData.VerifyTransaction(*tx) {
			info.Info = "交易验证失败"
			info.Status = false
			return info, errors.New("交易验证失败")
		}
		err = txpool.LocalTxPool.TxDataIN(*tx)
		if err != nil {
			info.Info = "交易入池失败"
			info.Status = false
			return info, err
		}
	}
	//elapsed := time.Since(start)
	//fmt.Println("PoolTxs have ", elapsed)
	//fmt.Println("ip ", util.LocalIP, "res send at ", time.UnixMicro(time.Now().UnixMicro()))
	return info, nil
}

// JoinCluster 加入集群
func (s *Service) JoinCluster(ctx context.Context, req *BcGrpc.ReqJoin) (info *BcGrpc.VerifyInfo, err error) {
	info = &BcGrpc.VerifyInfo{}
	info.Info = "可以加入集群"
	info.Status = true

	// 读取集群文件
	cluster, err := Cluster.LoadClusterFile("./ClusterInfo")
	if err != nil {
		info.Info = "读取集群文件失败"
		info.Status = false
		return info, err
	}
	fmt.Println(strconv.Itoa(3301) == req.LocalPort)
	// 判断节点是否在集群中
	for _, node := range cluster.Node {
		if node.IP == req.LocalPort && strconv.Itoa(node.Port) == req.LocalPort {
			info.Info = "节点已在集群中"
			info.Status = false
			return info, errors.New("节点已在集群中")
		}
	}
	// 判断加入密钥是否正确
	if !bytes.Equal(req.JoinKey, cluster.Key) {
		info.Info = "加入集群的密钥错误"
		info.Status = false
		return info, errors.New("密钥错误")
	}

	var node Cluster.Node
	node.IP = req.LocalIp
	node.Port, err = strconv.Atoi(req.LocalPort)
	if err != nil {
		log.Panic(err)
	}
	cluster.AddNodeToClusterFile(&node)

	return info, nil
}

// BroadcastNode 向集群中其他节点发送本节点的信息
func (s *Service) BroadcastNode(ctx context.Context, req *BcGrpc.NodeInfo) (info *BcGrpc.VerifyInfo, err error) {
	info = &BcGrpc.VerifyInfo{}
	info.Info = "发送信息成功"
	info.Status = true
	var node = &Cluster.Node{}
	node.IP = req.LocalIp
	port, err := strconv.Atoi(req.LocalPort)
	if err != nil {
		info.Info = "类型转换失败"
		info.Status = false
		return info, err
	}
	node.Port = port
	Cluster.LocalNode.AddNodeToClusterFile(node)
	return info, nil
}

func (s Service) Handle(ctx context.Context, data *BcGrpc.TypAndData) (info *BcGrpc.Info, err error) {
	info = &BcGrpc.Info{}
	// 进行处理
	//log.Println("rpc typ:", data.Typ)
	err = consensus.LocalPeer.Handle(int(data.Typ), data.Data)
	if err != nil {
		return nil, err
	}
	return info, nil
}

//func (s Service) Time(ctx context.Context, data *BcGrpc.Times) (info *BcGrpc.Info, err error) {
//	info = &BcGrpc.Info{}
//	// 进行处理
//	//log.Println("rpc typ:", data.Typ)
//	consensus.LocalPeer.SetTime(data.Time)
//	return info, nil
//}
