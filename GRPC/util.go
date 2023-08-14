package GRPC

import (
	"CoralDB/Cluster"
	BcGrpc "CoralDB/Proto/blockchain"
	BCData "CoralDB/blockchain/blockchain_data"
	BCTable "CoralDB/blockchain/blockchain_table"
	"CoralDB/blockqueue"
	"CoralDB/cache"
	"CoralDB/util"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strconv"
	"time"
)

var nodeIp []string

func BlockToGrpcDataBlock(block *BCData.Block) *BcGrpc.DataBlock {
	newGrpcBlock := &BcGrpc.DataBlock{
		ID:                int64(block.Round),
		CurrentBlockHash:  block.CurrentBlockHash,
		PreviousBlockHash: block.PreviousBlockHash,
		MerKelRoot:        block.MerKelRoot,
		TimeTamp:          block.TimeStamp,
	}
	newGrpcBlock.TxInfo = []*BcGrpc.DataTransaction{}
	for _, tx := range block.Transactions {
		newGrpcBlock.TxInfo = append(newGrpcBlock.TxInfo, &BcGrpc.DataTransaction{
			TxID:      tx.TxID,
			DataID:    tx.DataID,
			Table:     tx.Table,
			Key:       tx.Key,
			Value:     tx.Value,
			Possessor: tx.Possessor,
			TimeStamp: tx.TimeStamp,
			PublicKey: tx.PublicKey,
			Signature: tx.Signature,
		})
	}
	return newGrpcBlock
}

func BlockToGrpcTableBlock(block *BCTable.Block) *BcGrpc.TableBlock {
	newGrpcBlock := &BcGrpc.TableBlock{
		ID:                int64(block.ID),
		CurrentBlockHash:  block.CurrentBlockHash,
		PreviousBlockHash: block.PreviousBlockHash,
		MerKelRoot:        block.MerKelRoot,
		TimeTamp:          block.TimeStamp,
	}
	newGrpcBlock.TxInfo = []*BcGrpc.TableTransaction{}
	for _, tx := range block.Transactions {
		newGrpcBlock.TxInfo = append(newGrpcBlock.TxInfo, &BcGrpc.TableTransaction{
			TxID:             tx.TxID,
			Table:            tx.Table,
			PermissionTables: tx.PermissionTable,
			Possessor:        tx.Possessor,
			TimeStamp:        tx.TimeStamp,
			PublicKey:        tx.PublicKey,
			Signature:        tx.Signature,
		})
	}
	return newGrpcBlock
}

func GrpcDataBlockToBlock(block *BcGrpc.DataBlock) *BCData.Block {
	newBlock := &BCData.Block{
		Round:             uint64(block.ID),
		CurrentBlockHash:  block.CurrentBlockHash,
		PreviousBlockHash: block.PreviousBlockHash,
		MerKelRoot:        block.MerKelRoot,
		TimeStamp:         block.TimeTamp,
	}
	newBlock.Transactions = []*BCData.Transaction{}
	for _, tx := range block.TxInfo {
		newBlock.Transactions = append(newBlock.Transactions, &BCData.Transaction{
			TxID:      tx.TxID,
			DataID:    tx.DataID,
			Table:     tx.Table,
			Key:       tx.Key,
			Value:     tx.Value,
			Possessor: tx.Possessor,
			TimeStamp: tx.TimeStamp,
			PublicKey: tx.PublicKey,
			Signature: tx.Signature,
		})
	}
	return newBlock
}

func GrpcTableBlockToBlock(block *BcGrpc.TableBlock) *BCTable.Block {
	newBlock := &BCTable.Block{
		ID:                int(block.ID),
		CurrentBlockHash:  block.CurrentBlockHash,
		PreviousBlockHash: block.PreviousBlockHash,
		MerKelRoot:        block.MerKelRoot,
		TimeStamp:         block.TimeTamp,
	}
	newBlock.Transactions = []*BCTable.Transaction{}
	for _, tx := range block.TxInfo {
		newBlock.Transactions = append(newBlock.Transactions, &BCTable.Transaction{
			TxID:            tx.TxID,
			Table:           tx.Table,
			PermissionTable: tx.PermissionTables,
			Possessor:       tx.Possessor,
			TimeStamp:       tx.TimeStamp,
			PublicKey:       tx.PublicKey,
			Signature:       tx.Signature,
		})
	}
	return newBlock
}

// GrpcDataTxToDataTx 交易的转换
func GrpcDataTxToDataTx(tx *BcGrpc.DataTransaction) *BCData.Transaction {
	newTx := &BCData.Transaction{
		TxID:      tx.TxID,
		DataID:    tx.DataID,
		Table:     tx.Table,
		Key:       tx.Key,
		Value:     tx.Value,
		Possessor: tx.Possessor,
		TimeStamp: tx.TimeStamp,
		PublicKey: tx.PublicKey,
		Signature: tx.Signature,
	}
	return newTx
}

func GrpcTableTxToTableTx(tx *BcGrpc.TableTransaction) *BCTable.Transaction {
	newTx := &BCTable.Transaction{
		TxID:            tx.TxID,
		Table:           tx.Table,
		PermissionTable: tx.PermissionTables,
		Possessor:       tx.Possessor,
		TimeStamp:       tx.TimeStamp,
		PublicKey:       tx.PublicKey,
		Signature:       tx.Signature,
	}
	return newTx
}

func DataTxToGrpcDataTx(tx *BCData.Transaction) *BcGrpc.DataTransaction {
	newTx := &BcGrpc.DataTransaction{
		TxID:      tx.TxID,
		DataID:    tx.DataID,
		Table:     tx.Table,
		Key:       tx.Key,
		Value:     tx.Value,
		Possessor: tx.Possessor,
		TimeStamp: tx.TimeStamp,
		PublicKey: tx.PublicKey,
		Signature: tx.Signature,
	}
	return newTx
}

func TableTxToGrpcTableTx(tx *BCTable.Transaction) *BcGrpc.TableTransaction {
	newTx := &BcGrpc.TableTransaction{
		TxID:             tx.TxID,
		Table:            tx.Table,
		PermissionTables: tx.PermissionTable,
		Possessor:        tx.Possessor,
		TimeStamp:        tx.TimeStamp,
		PublicKey:        tx.PublicKey,
		Signature:        tx.Signature,
	}
	return newTx
}

// BroadCast 广播自己已加入集群
func BroadCast(ip string, port int) {
	// 参数为本节点向哪个节点提交的申请
	length := len(Cluster.LocalNode.Node)
	for i := 0; i < length; i++ {
		node := Cluster.LocalNode.Node[i]
		if (node.IP == util.LocalIP && node.Port == util.LocalPort) || (node.IP == ip && node.Port == port) {
			continue
		}
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.IP, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("%s:%d 网络异常", node.IP, node.Port)
		}

		// 获得Grpc句柄
		c := BcGrpc.NewBlockChainServiceClient(conn)
		// 通过句柄调用函数
		_, err = c.BroadcastNode(context.Background(), &BcGrpc.NodeInfo{
			LocalIp:   node.IP,
			LocalPort: strconv.Itoa(node.Port),
		})
		if err != nil {
			log.Panic(err)
		}
		conn.Close()
	}
}

// JoinToCluster 加入集群
func JoinToCluster(ip string, port int, key string) error {
	var node Cluster.Node
	fmt.Println(ip)
	// 获得本机的IP和端口号
	node.IP = util.GetLocalIp()
	node.Port = port
	// 获得密钥的hash
	KeyHash := sha256.Sum256([]byte(key))
	Key := KeyHash[:]
	// 提交加入集群的请求(向已加入集群的节点)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Network exception!")
		log.Panic(err)
	}
	// 获得grpc句柄
	client := BcGrpc.NewBlockChainServiceClient(conn)
	// 通过句柄调用函数，验证节点是否可以加入集群
	_, err = client.JoinCluster(context.Background(), &BcGrpc.ReqJoin{
		LocalIp:   node.IP,
		LocalPort: strconv.Itoa(node.Port),
		JoinKey:   Key,
	})
	if err != nil {
		log.Panic(err)
	}
	// 判断集群文件是否存在
	_, err = os.Stat("./ClusterInfo")
	if !os.IsNotExist(err) { // 如果没有报错即集群文件存在
		// 如果存在就先清除集群文件
		err = os.Remove("./ClusterInfo")
		if err != nil {
			log.Panic(err)
		}
	}
	// 接收集群文件
	flag := Cluster.Client(ip)
	if !flag {
		return errors.New("接收集群文件失败")
	}
	// 读取集群文件
	clu, err := Cluster.LoadClusterFile("./ClusterInfo")
	if err != nil {
		log.Panic(err)
	}
	// 将自己加入集群
	//clu.AddNodeToClusterFile(&node)
	// 全局变量的覆盖
	Cluster.LocalNode = clu
	err = conn.Close()
	if err != nil {
		log.Panic(err)
	}
	return nil
}

// SubmitDataTransaction 交易提交到交易池
func SubmitDataTransaction(tx BCData.Transaction) {
	num := len(nodeIp)
	switch num {
	case 3:
		go handle(nodeIp[0], tx)
		go handle(nodeIp[1], tx)
		go handle(nodeIp[2], tx)
	case 7:
		go handle(nodeIp[0], tx)
		go handle(nodeIp[1], tx)
		go handle(nodeIp[2], tx)
		go handle(nodeIp[3], tx)
		go handle(nodeIp[4], tx)
		go handle(nodeIp[5], tx)
		go handle(nodeIp[6], tx)
	case 11:
		go handle(nodeIp[0], tx)
		go handle(nodeIp[1], tx)
		go handle(nodeIp[2], tx)
		go handle(nodeIp[3], tx)
		go handle(nodeIp[4], tx)
		go handle(nodeIp[5], tx)
		go handle(nodeIp[6], tx)
		go handle(nodeIp[7], tx)
		go handle(nodeIp[8], tx)
		go handle(nodeIp[9], tx)
		go handle(nodeIp[10], tx)
	case 15:
		go handle(nodeIp[0], tx)
		go handle(nodeIp[1], tx)
		go handle(nodeIp[2], tx)
		go handle(nodeIp[3], tx)
		go handle(nodeIp[4], tx)
		go handle(nodeIp[5], tx)
		go handle(nodeIp[6], tx)
		go handle(nodeIp[7], tx)
		go handle(nodeIp[8], tx)
		go handle(nodeIp[9], tx)
		go handle(nodeIp[10], tx)
		go handle(nodeIp[11], tx)
		go handle(nodeIp[12], tx)
		go handle(nodeIp[13], tx)
		go handle(nodeIp[14], tx)
	case 19:
		go handle(nodeIp[0], tx)
		go handle(nodeIp[1], tx)
		go handle(nodeIp[2], tx)
		go handle(nodeIp[3], tx)
		go handle(nodeIp[4], tx)
		go handle(nodeIp[5], tx)
		go handle(nodeIp[6], tx)
		go handle(nodeIp[7], tx)
		go handle(nodeIp[8], tx)
		go handle(nodeIp[9], tx)
		go handle(nodeIp[10], tx)
		go handle(nodeIp[11], tx)
		go handle(nodeIp[12], tx)
		go handle(nodeIp[13], tx)
		go handle(nodeIp[14], tx)
		go handle(nodeIp[15], tx)
		go handle(nodeIp[16], tx)
		go handle(nodeIp[17], tx)
		go handle(nodeIp[18], tx)
	case 24:
		go handle(nodeIp[0], tx)
		go handle(nodeIp[1], tx)
		go handle(nodeIp[2], tx)
		go handle(nodeIp[3], tx)
		go handle(nodeIp[4], tx)
		go handle(nodeIp[5], tx)
		go handle(nodeIp[6], tx)
		go handle(nodeIp[7], tx)
		go handle(nodeIp[8], tx)
		go handle(nodeIp[9], tx)
		go handle(nodeIp[10], tx)
		go handle(nodeIp[11], tx)
		go handle(nodeIp[12], tx)
		go handle(nodeIp[13], tx)
		go handle(nodeIp[14], tx)
		go handle(nodeIp[15], tx)
		go handle(nodeIp[16], tx)
		go handle(nodeIp[17], tx)
		go handle(nodeIp[18], tx)
		go handle(nodeIp[19], tx)
		go handle(nodeIp[20], tx)
		go handle(nodeIp[21], tx)
		go handle(nodeIp[22], tx)
		go handle(nodeIp[23], tx)
	}

	// 向所有节点发送交易信息
	//for _, node := range Cluster.LocalNode.Node {
	//	if node.IP == util.LocalIP && node.Port == util.LocalPort {
	//		continue
	//	}
	//	//go func(ip string, port int) {
	//	//start := time.Now() // 获取当前时间
	//	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.IP, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	//	defer conn.Close()
	//	if err != nil {
	//		fmt.Println("Network exception!")
	//		log.Panic(err)
	//	}
	//	// 获得grpc句柄
	//	client := BcGrpc.NewBlockChainServiceClient(conn)
	//	// 通过句柄调用函数
	//	_, err = client.DataTradingPool(context.Background(), &BcGrpc.DataTransaction{
	//		TxID:      tx.TxID,
	//		DataID:    tx.DataID,
	//		Table:     tx.Table,
	//		Key:       tx.Key,
	//		Value:     tx.Value,
	//		Possessor: tx.Possessor,
	//		TimeStamp: tx.TimeStamp,
	//		PublicKey: tx.PublicKey,
	//		Signature: tx.Signature,
	//	})
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//	err = conn.Close()
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//elapsed := time.Since(start)
	//fmt.Println("broadcast time ", elapsed)
	//}(ip, port)

	//}
}

func handle(addr string, tx BCData.Transaction) {
	//start := time.Now() // 获取当前时间
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		fmt.Println("Network exception!")
		log.Panic(err)
	}
	// 获得grpc句柄
	client := BcGrpc.NewBlockChainServiceClient(conn)
	// 通过句柄调用函数
	_, err = client.DataTradingPool(context.Background(), &BcGrpc.DataTransaction{
		TxID:      tx.TxID,
		DataID:    tx.DataID,
		Table:     tx.Table,
		Key:       tx.Key,
		Value:     tx.Value,
		Possessor: tx.Possessor,
		TimeStamp: tx.TimeStamp,
		PublicKey: tx.PublicKey,
		Signature: tx.Signature,
	})
	if err != nil {
		log.Panic(err)
	}
	err = conn.Close()
	if err != nil {
		log.Panic(err)
	}
	//elapsed := time.Since(start)
	//fmt.Printf("%s send tx %d use ", addr, num)
	//fmt.Println(elapsed)
}

// SubmitTableTransaction 交易的提交到交易池
func SubmitTableTransaction(tx BCTable.Transaction) {
	// 向集群节点发送交易信息
	for _, node := range Cluster.LocalNode.Node {
		if node.IP == util.LocalIP && node.Port == util.LocalPort {
			continue
		}
		nodeIp = append(nodeIp, fmt.Sprintf("%s:%d", node.IP, node.Port))
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.IP, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer conn.Close()
		if err != nil {
			fmt.Println("Network exception!")
			log.Panic(err)
		}
		// 获得grpc句柄
		client := BcGrpc.NewBlockChainServiceClient(conn)
		// 通过句柄调用函数
		_, err = client.TableTradingPool(context.Background(), &BcGrpc.TableTransaction{
			TxID:             tx.TxID,
			Table:            tx.Table,
			PermissionTables: tx.PermissionTable,
			Possessor:        tx.Possessor,
			TimeStamp:        tx.TimeStamp,
			PublicKey:        tx.PublicKey,
			Signature:        tx.Signature,
		})
		if err != nil {
			log.Panic(err)
		}
		err = conn.Close()
		if err != nil {
			log.Panic(err)
		}
	}
}

var nums = 0

// SubmitDataTransactions 多笔交易的入池
func SubmitDataTransactions() {
	for {
		time.Sleep(time.Millisecond * 10)
		if blockqueue.LocalDataTxQueue.Size >= 1 {
			var txs []BCData.Transaction
			for i := 0; i < blockqueue.LocalDataTxQueue.Size; i++ {
				get := blockqueue.LocalDataTxQueue.Get()
				tx := get.(BCData.Transaction)
				txs = append(txs, tx)
			}
			num := len(nodeIp)
			switch num {
			case 3:
				go handleTxs(nodeIp[0], txs, nums)
				go handleTxs(nodeIp[1], txs, nums)
				go handleTxs(nodeIp[2], txs, nums)
			case 7:
				go handleTxs(nodeIp[0], txs, nums)
				go handleTxs(nodeIp[1], txs, nums)
				go handleTxs(nodeIp[2], txs, nums)
				go handleTxs(nodeIp[3], txs, nums)
				go handleTxs(nodeIp[4], txs, nums)
				go handleTxs(nodeIp[5], txs, nums)
				go handleTxs(nodeIp[6], txs, nums)
			case 11:
				go handleTxs(nodeIp[0], txs, nums)
				go handleTxs(nodeIp[1], txs, nums)
				go handleTxs(nodeIp[2], txs, nums)
				go handleTxs(nodeIp[3], txs, nums)
				go handleTxs(nodeIp[4], txs, nums)
				go handleTxs(nodeIp[5], txs, nums)
				go handleTxs(nodeIp[6], txs, nums)
				go handleTxs(nodeIp[7], txs, nums)
				go handleTxs(nodeIp[8], txs, nums)
				go handleTxs(nodeIp[9], txs, nums)
				go handleTxs(nodeIp[10], txs, nums)
			case 15:
				go handleTxs(nodeIp[0], txs, nums)
				go handleTxs(nodeIp[1], txs, nums)
				go handleTxs(nodeIp[2], txs, nums)
				go handleTxs(nodeIp[3], txs, nums)
				go handleTxs(nodeIp[4], txs, nums)
				go handleTxs(nodeIp[5], txs, nums)
				go handleTxs(nodeIp[6], txs, nums)
				go handleTxs(nodeIp[7], txs, nums)
				go handleTxs(nodeIp[8], txs, nums)
				go handleTxs(nodeIp[9], txs, nums)
				go handleTxs(nodeIp[10], txs, nums)
				go handleTxs(nodeIp[11], txs, nums)
				go handleTxs(nodeIp[12], txs, nums)
				go handleTxs(nodeIp[13], txs, nums)
				go handleTxs(nodeIp[14], txs, nums)
			case 19:
				go handleTxs(nodeIp[0], txs, nums)
				go handleTxs(nodeIp[1], txs, nums)
				go handleTxs(nodeIp[2], txs, nums)
				go handleTxs(nodeIp[3], txs, nums)
				go handleTxs(nodeIp[4], txs, nums)
				go handleTxs(nodeIp[5], txs, nums)
				go handleTxs(nodeIp[6], txs, nums)
				go handleTxs(nodeIp[7], txs, nums)
				go handleTxs(nodeIp[8], txs, nums)
				go handleTxs(nodeIp[9], txs, nums)
				go handleTxs(nodeIp[10], txs, nums)
				go handleTxs(nodeIp[11], txs, nums)
				go handleTxs(nodeIp[12], txs, nums)
				go handleTxs(nodeIp[13], txs, nums)
				go handleTxs(nodeIp[14], txs, nums)
				go handleTxs(nodeIp[15], txs, nums)
				go handleTxs(nodeIp[16], txs, nums)
				go handleTxs(nodeIp[17], txs, nums)
				go handleTxs(nodeIp[18], txs, nums)
			case 24:
				go handleTxs(nodeIp[0], txs, nums)
				go handleTxs(nodeIp[1], txs, nums)
				go handleTxs(nodeIp[2], txs, nums)
				go handleTxs(nodeIp[3], txs, nums)
				go handleTxs(nodeIp[4], txs, nums)
				go handleTxs(nodeIp[5], txs, nums)
				go handleTxs(nodeIp[6], txs, nums)
				go handleTxs(nodeIp[7], txs, nums)
				go handleTxs(nodeIp[8], txs, nums)
				go handleTxs(nodeIp[9], txs, nums)
				go handleTxs(nodeIp[10], txs, nums)
				go handleTxs(nodeIp[11], txs, nums)
				go handleTxs(nodeIp[12], txs, nums)
				go handleTxs(nodeIp[13], txs, nums)
				go handleTxs(nodeIp[14], txs, nums)
				go handleTxs(nodeIp[15], txs, nums)
				go handleTxs(nodeIp[16], txs, nums)
				go handleTxs(nodeIp[17], txs, nums)
				go handleTxs(nodeIp[18], txs, nums)
				go handleTxs(nodeIp[19], txs, nums)
				go handleTxs(nodeIp[20], txs, nums)
				go handleTxs(nodeIp[21], txs, nums)
				go handleTxs(nodeIp[22], txs, nums)
				go handleTxs(nodeIp[23], txs, nums)
			}
			nums++
		}
	}
}

// 多笔交易的处理
func handleTxs(addr string, txs []BCData.Transaction, nums int) {
	//start := time.Now()
	//log.Println("addr", addr, "time", time.Now().UnixMicro(), "round", nums)
	//fmt.Println("addr", addr, " req send time", time.UnixMicro(time.Now().UnixMicro()), "round", nums)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		fmt.Printf("%s 网络异常", addr)
	}
	var GrpcTxs BcGrpc.DataTransactions
	for i := 0; i < len(txs); i++ {
		GrpcTxs.Transactions = append(GrpcTxs.Transactions, DataTxToGrpcDataTx(&txs[i]))
	}
	// 获得grpc句柄
	client := BcGrpc.NewBlockChainServiceClient(conn)
	// 调用函数

	_, err = client.DataTradingPoolTxs(context.Background(), &GrpcTxs)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println("addr", addr, " res in time", time.UnixMicro(time.Now().UnixMicro()), "round", nums)
	err = conn.Close()
	if err != nil {
		log.Panic(err)
	}
	//elapsed := time.Since(start)
	//fmt.Println("addr", addr, "sendTx have ", elapsed, "round", nums)
}

// DataBlockDistribute 数据区块的分发
func DataBlockDistribute() {
	for {
		time.Sleep(time.Microsecond * 100)
		if blockqueue.LocalDataBlockQueue.Size >= 1 {
			//log.Println("RPC数据区块区块的分发")
			get := blockqueue.LocalDataBlockQueue.Get()
			block := get.(*BCData.Block)
			num := len(nodeIp)
			switch num {
			case 1:
				go handleBlock(nodeIp[0], block)
			case 2:
				go handleBlock(nodeIp[0], block)
				go handleBlock(nodeIp[1], block)
			case 3:
				go handleBlock(nodeIp[0], block)
				go handleBlock(nodeIp[1], block)
				go handleBlock(nodeIp[2], block)
			case 4:
				go handleBlock(nodeIp[0], block)
				go handleBlock(nodeIp[1], block)
				go handleBlock(nodeIp[2], block)
				go handleBlock(nodeIp[3], block)
			case 7:
				go handleBlock(nodeIp[0], block)
				go handleBlock(nodeIp[1], block)
				go handleBlock(nodeIp[2], block)
				go handleBlock(nodeIp[3], block)
				go handleBlock(nodeIp[4], block)
				go handleBlock(nodeIp[5], block)
				go handleBlock(nodeIp[6], block)

			case 11:
				go handleBlock(nodeIp[0], block)
				go handleBlock(nodeIp[1], block)
				go handleBlock(nodeIp[2], block)
				go handleBlock(nodeIp[3], block)
				go handleBlock(nodeIp[4], block)
				go handleBlock(nodeIp[5], block)
				go handleBlock(nodeIp[6], block)
				go handleBlock(nodeIp[7], block)
				go handleBlock(nodeIp[8], block)
				go handleBlock(nodeIp[9], block)
				go handleBlock(nodeIp[10], block)
			case 15:
				go handleBlock(nodeIp[0], block)
				go handleBlock(nodeIp[1], block)
				go handleBlock(nodeIp[2], block)
				go handleBlock(nodeIp[3], block)
				go handleBlock(nodeIp[4], block)
				go handleBlock(nodeIp[5], block)
				go handleBlock(nodeIp[6], block)
				go handleBlock(nodeIp[7], block)
				go handleBlock(nodeIp[8], block)
				go handleBlock(nodeIp[9], block)
				go handleBlock(nodeIp[10], block)
				go handleBlock(nodeIp[11], block)
				go handleBlock(nodeIp[12], block)
				go handleBlock(nodeIp[13], block)
				go handleBlock(nodeIp[14], block)
			case 19:
				go handleBlock(nodeIp[0], block)
				go handleBlock(nodeIp[1], block)
				go handleBlock(nodeIp[2], block)
				go handleBlock(nodeIp[3], block)
				go handleBlock(nodeIp[4], block)
				go handleBlock(nodeIp[5], block)
				go handleBlock(nodeIp[6], block)
				go handleBlock(nodeIp[7], block)
				go handleBlock(nodeIp[8], block)
				go handleBlock(nodeIp[9], block)
				go handleBlock(nodeIp[10], block)
				go handleBlock(nodeIp[11], block)
				go handleBlock(nodeIp[12], block)
				go handleBlock(nodeIp[13], block)
				go handleBlock(nodeIp[14], block)
				go handleBlock(nodeIp[15], block)
				go handleBlock(nodeIp[16], block)
				go handleBlock(nodeIp[17], block)
				go handleBlock(nodeIp[18], block)
			case 24:
				go handleBlock(nodeIp[0], block)
				go handleBlock(nodeIp[1], block)
				go handleBlock(nodeIp[2], block)
				go handleBlock(nodeIp[3], block)
				go handleBlock(nodeIp[4], block)
				go handleBlock(nodeIp[5], block)
				go handleBlock(nodeIp[6], block)
				go handleBlock(nodeIp[7], block)
				go handleBlock(nodeIp[8], block)
				go handleBlock(nodeIp[9], block)
				go handleBlock(nodeIp[10], block)
				go handleBlock(nodeIp[11], block)
				go handleBlock(nodeIp[12], block)
				go handleBlock(nodeIp[13], block)
				go handleBlock(nodeIp[14], block)
				go handleBlock(nodeIp[15], block)
				go handleBlock(nodeIp[16], block)
				go handleBlock(nodeIp[17], block)
				go handleBlock(nodeIp[18], block)
				go handleBlock(nodeIp[19], block)
				go handleBlock(nodeIp[20], block)
				go handleBlock(nodeIp[21], block)
				go handleBlock(nodeIp[22], block)
				go handleBlock(nodeIp[23], block)
			}

			//for _, node := range Cluster.LocalNode.Node {
			//	if node.IP == util.LocalIP && node.Port == util.LocalPort {
			//		continue
			//	}
			//	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.IP, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
			//	defer conn.Close()
			//	if err != nil {
			//		fmt.Printf("%s:%d 网络异常", node.IP, node.Port)
			//	}
			//	// 获得grpc句柄
			//	client := BcGrpc.NewBlockChainServiceClient(conn)
			//	// 调用函数
			//	_, err = client.DistributeDataBlock(context.Background(), BlockToGrpcDataBlock(block))
			//	if err != nil {
			//		log.Panic(err)
			//	}
			//	err = conn.Close()
			//	if err != nil {
			//		log.Panic(err)
			//	}
			//}
		}
	}
}

func handleBlock(addr string, block *BCData.Block) {
	//start := time.Now() // 获取当前时间
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		fmt.Printf("%s 网络异常", addr)
	}
	// 获得grpc句柄
	client := BcGrpc.NewBlockChainServiceClient(conn)
	// 调用函数
	_, err = client.DistributeDataBlock(context.Background(), BlockToGrpcDataBlock(block))
	if err != nil {
		log.Println(err)
	}
	err = conn.Close()
	if err != nil {
		log.Panic(err)
	}
	//elapsed := time.Since(start)
	//fmt.Printf("%s send block use ", addr)
	//fmt.Println(elapsed)
}

// TableBlockDistribute 权限区块的分发
func TableBlockDistribute() {
	for {
		time.Sleep(time.Millisecond * 500)
		if blockqueue.LocalTableBlockQueue.Size >= 1 {
			//fmt.Println("RPC权限区块区块的分发")
			get := blockqueue.LocalTableBlockQueue.Get()
			block := get.(BCTable.Block)
			for _, node := range Cluster.LocalNode.Node {
				if node.IP == util.LocalIP && node.Port == util.LocalPort {
					continue
				}
				conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.IP, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
				defer conn.Close()
				if err != nil {
					fmt.Printf("%s:%d 网络异常", node.IP, node.Port)
				}
				// 获得grpc句柄
				client := BcGrpc.NewBlockChainServiceClient(conn)
				// 调用函数
				_, err = client.DistributeTableBlock(context.Background(), BlockToGrpcTableBlock(&block))
				if err != nil {
					log.Panic(err)
				}
				//log.Println("表区块的分发")
				err = conn.Close()
				if err != nil {
					log.Panic(err)
				}
			}

		}
	}

}

// DataBlockSynchronization 数据区块的同步
func DataBlockSynchronization() {
	time.Sleep(time.Second * 2)
	// 获得本地最后区块的信息
	lastHash := BCData.LocalDataBlockChain.TailHash
	block, err := BCData.LocalDataBlockChain.GetBlockByHash(lastHash)
	if err != nil {
		log.Panic(err)
	}
	// 得到记账节点的IP
	for _, node := range Cluster.LocalNode.Node {
		if node.Account == true {
			// 向记账节点发出申请
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.IP, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer conn.Close()
			if err != nil {
				fmt.Printf("%s:%d 网络异常", node.IP, node.Port)
			}
			// 获得grpc句柄
			client := BcGrpc.NewBlockChainServiceClient(conn)
			// 通过句柄调用函数
			re, err := client.DataBlockSynchronization(context.Background(), &BcGrpc.ReqDataBlock{
				BlockID: int64(block.Round),
				Hash:    block.CurrentBlockHash,
			})
			// 区块上链
			for i := 0; i < len(re.Blocks); i++ {
				// 区块的转换
				newBlock := GrpcDataBlockToBlock(re.Blocks[len(re.Blocks)-1-i])
				//上链
				BCData.LocalDataBlockChain.AddBlockToChain(*newBlock)
				// 缓存的更新
				cache.LocalCache.UpdateByDataBlock(*newBlock)
			}
		}
	}
	// 输出信息
	fmt.Println("数据区块同步完成")

}

// TableBlockSynchronization 权限区块的同步
func TableBlockSynchronization() {
	time.Sleep(time.Second * 2)
	// 获得本地最后区块的信息
	lastHash := BCTable.LocalTableBlockChain.TailHash
	block, err := BCTable.LocalTableBlockChain.GetBlockByHash(lastHash)
	if err != nil {
		log.Panic(err)
	}
	// 得到记账节点的IP
	for _, node := range Cluster.LocalNode.Node {
		if node.Account == true {
			// 向记账节点发出申请
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.IP, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer conn.Close()
			if err != nil {
				fmt.Printf("%s:%d 网络异常", node.IP, node.Port)
			}
			// 获得grpc句柄
			client := BcGrpc.NewBlockChainServiceClient(conn)
			// 通过句柄调用函数
			re, err := client.TableBlockSynchronization(context.Background(), &BcGrpc.ReqTableBlock{
				BlockID: int64(block.ID),
				Hash:    block.CurrentBlockHash,
			})
			// 区块上链
			for i := 0; i < len(re.Blocks); i++ {
				// 区块的转换
				newBlock := GrpcTableBlockToBlock(re.Blocks[len(re.Blocks)-1-i])
				//上链
				BCTable.LocalTableBlockChain.AddBlockToChain(*newBlock)
				// 缓存的更新
				cache.LocalCache.UpdateByTableBlock(*newBlock)
			}

		}
	}
	fmt.Println("权限区块同步完成")
}
