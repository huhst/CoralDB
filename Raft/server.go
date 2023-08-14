package Raft

import (
	"CoralDB/Cluster"
	"CoralDB/txpool"
	"CoralDB/util"
	"strconv"
)

//定义节点数量
var nodeCount = 3

//节点池
var nodePool int

//选举超时时间（秒）
var electionTimeout = 3

//心跳检测超时时间
var heartBeatTimeout = 4

//心跳检测频率（秒）
var heartBeatRate = 2

//// MessageStore 存储信息
//var MessageStore = make(map[int]string)

func Start(tpl *txpool.TxPool) {
	nodeCount = len(Cluster.LocalNode.Node)
	nodePool = 3302

	//传入节点编号，端口号，创建raft实例
	raft := NewRaft(util.LocalIP, strconv.Itoa(nodePool))

	//注册rpc服务绑定http协议上开启监听
	go rpcRegister(raft)

	//发送心跳包
	go raft.sendHeartPacket()

	//尝试成为候选人并选举
	go raft.tryToBeCandidateWithElection()

	//进行心跳超时检测
	go raft.heartTimeoutDetection()

	// 记账权的更新
	go raft.SetAccount(tpl)
}
