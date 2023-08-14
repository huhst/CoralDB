package Cluster

import (
	//"blockchainTest/util"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// SendFile 发送文件
func SendFile(conn net.Conn, filePath string) {
	//只读打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	buf := make([]byte, 4096)
	for {
		//从本地文件中读数据，写给网络接收端。读多少，写多少
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("\n发送文件完毕\n")
				fmt.Printf("> ")
				return
			} else {
				fmt.Printf("file.Read()方法执行错误,错误为:%v\n", err)
				return
			}
		}
		//写到网络socket中
		_, err = conn.Write(buf[:n])
	}
}

// RecvFiles 接收文件
func RecvFiles(conn net.Conn, fileName string) {
	//按照文件名创建新文件
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("os.Create()函数执行错误，错误为:%v\n", err)
		return
	}
	defer file.Close()

	//从网络中读数据，写入本地文件
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)

		//写入本地文件，读多少，写多少
		file.Write(buf[:n])
		if err != nil {
			if err == io.EOF {
				//util.FlagJoin = true
				fmt.Printf("接收文件完成\n")
			} else {
				fmt.Printf("conn.Read()方法执行出错，错误为:%v\n", err)
			}
			return
		}
	}
}
