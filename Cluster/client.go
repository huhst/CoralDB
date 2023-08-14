package Cluster

import (
	"fmt"
	"log"
	"net"
)

// Client 发送加入集群的请求与接受Cluster文件
func Client(ServerIp string) bool {
	address := fmt.Sprintf("%s:%d", ServerIp, 3303)
	conn, err := net.Dial("tcp", address)
	// 关闭连接
	//defer conn.Close()
	if err != nil {
		log.Panic(err)
	}
	// 判断本地是否有集群文件，有就删除
	RecvFiles(conn, "./ClusterInfo") //接收文件
	return true
}
