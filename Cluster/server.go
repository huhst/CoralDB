package Cluster

import (
	"log"
	"net"
)

// 监听是否有节点想要加入集群与genesis文件的发送

// Server 启动监听服务
func Server() {
	// 启动监听程序
	listener, err := net.Listen("tcp", ":3303")
	//defer listener.Close()
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := listener.Accept()
		//defer conn.Close()
		if err != nil {
			log.Panic(err)
		}
		Handler(conn)
	}

}

// Handler 处理函数
func Handler(conn net.Conn) {
	// 关闭连接
	defer conn.Close()
	// 读取genesis文件与文件的发送
	SendFile(conn, "./ClusterInfo")
}
