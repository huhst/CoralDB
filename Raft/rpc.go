package Raft

import (
	"CoralDB/Cluster"
	"CoralDB/txpool"
	"CoralDB/util"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"
)

//注册rpc服务绑定http协议上开启监听
func rpcRegister(raft *Raft) {
	//注册一个RPC服务器
	if err := rpc.Register(raft); err != nil {
		log.Panicln("注册RPC失败", err)
	}
	//ip := raft.node.IP
	//port := raft.node.Port
	//把RPC服务绑定到http协议上c
	rpc.HandleHTTP()
	//127.0.0.1:6870|6871|6872
	err := http.ListenAndServe(":3302", nil)
	if err != nil {
		log.Panic(err)
	}
}

//广播，调用所有节点的method方法（不广播自己）
func (rf *Raft) broadcast(method string, args interface{}, fun func(ok bool)) {
	for _, node := range Cluster.LocalNode.Node {
		//不广播自己
		if node.IP+"3302" == rf.me {
			continue
		}
		//连接远程节点的rpc
		conn, err := rpc.DialHTTP("tcp", node.IP+":3302")
		if err != nil {
			//连接失败，调用回调
			fun(false)
			continue
		}
		var bo bool
		err = conn.Call(method, args, &bo)
		if err != nil {
			//调用失败，调用回调
			fun(false)
			continue
		}
		//回调
		fun(bo)
		err = conn.Close()
		if err != nil {
			log.Panic(err)
		}

	}
}

// HeartBeatResponse 心跳检测回复
func (rf *Raft) HeartBeatResponse(node NodeInfo, b *bool) error {
	//因为发送心跳的一定是leader，之所以写这一句的目的是如果有down的节点恢复了，直接是follower，所以直接告诉它leader是谁即可
	rf.setCurrentLeader(node.IP)
	//最后一次心跳的时间
	rf.lastSendHeartBeatTime = millisecond()
	//fmt.Printf("收到来自leader[%s]节点的心跳检测\n", node.IP)
	*b = true
	return nil
}

// ConfirmationLeader 确认领导者
func (rf *Raft) ConfirmationLeader(node NodeInfo, b *bool) error {
	rf.setCurrentLeader(node.IP)
	*b = true
	//fmt.Println(node.IP, "成为了领导者")
	fmt.Printf("> ")
	rf.reDefault()
	// 将记账权给领导者
	for _, nodes := range Cluster.LocalNode.Node {
		if nodes.IP == node.IP {
			nodes.Account = true
		} else {
			nodes.Account = false
		}
	}
	// 文件更新
	Cluster.SaveClusterFile()
	return nil
}

// Vote 投票
func (rf *Raft) Vote(node NodeInfo, b *bool) error {
	if rf.votedFor == "-1" && rf.currentLeader == "-1" {
		rf.setVoteFor(node.IP)
		//fmt.Printf("投票成功，已投%s节点\n", node.IP)
		*b = true
	} else {
		*b = false
	}
	return nil
}

// SetAccount 记账权的更新
func (rf *Raft) SetAccount(tpl *txpool.TxPool) {
	for {
		time.Sleep(time.Millisecond * 500)
		if rf.state == 2 {
			for _, node := range Cluster.LocalNode.Node {
				if node.IP == util.LocalIP && node.Port == util.LocalPort {
					node.Account = true
					util.LocalIsAccount = true
					tpl.SetMod(util.LocalIsAccount)
				}
			}
		} else {
			for _, node := range Cluster.LocalNode.Node {
				if node.IP == util.LocalIP && node.Port == util.LocalPort {
					node.Account = false
					util.LocalIsAccount = false
					tpl.SetMod(util.LocalIsAccount)
				}
			}
		}
	}

}
