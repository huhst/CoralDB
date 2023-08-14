package Cluster

import (
	"CoralDB/util"
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Node  定义集群节点结构
type Node struct {
	IP      string // IP
	Port    int    // 端口号
	Account bool   //是否有记账权
}

// Cluster 定义集群的结构
type Cluster struct {
	Node []*Node // 集群中的节点
	Key  []byte  // 加入集群的密钥
}

var LocalNode *Cluster

// Init 集群的初始化
func Init(key string) *Cluster {
	// 设置密钥
	Key := sha256.Sum256([]byte(key))
	// 节点信息
	var nodes []*Node
	var node Node
	node.IP = util.GetLocalIp()
	node.Port = 3301
	node.Account = false
	nodes = append(nodes, &node)
	clu := Cluster{
		Node: nodes,
		Key:  Key[:],
	}
	LocalNode = &clu
	return &clu
}

// SaveClusterFile 集群写入文件
func SaveClusterFile() {
	var buf bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(LocalNode)
	if err != nil {
		fmt.Println("SaveClusterFile failed")
		log.Panic(err)
	}
	err = ioutil.WriteFile("ClusterInfo", util.AesCTREncrypt(buf.Bytes(), []byte("1234567812345678")), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// LoadClusterFile 读取集群文件内容
func LoadClusterFile(fileName string) (*Cluster, error) {
	var clu = &Cluster{}
	// 判断文件是否存在
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, err
	}
	// 读取文件
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	// 解码
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(util.AesCTRDecrypt(content, []byte("1234567812345678"))))
	err = decoder.Decode(&clu)
	if err != nil {
		log.Panic(err)
	}
	LocalNode = clu

	return LocalNode, nil
}

// AddNodeToClusterFile 向集群文件添加节点信息
func (clu *Cluster) AddNodeToClusterFile(node *Node) {
	clu.Node = append(clu.Node, node)
	LocalNode = clu
	SaveClusterFile() // 将文件进行保存
}

// UpdateClusterFile 集群文件的更新（删除重复的节点）
func (clu *Cluster) UpdateClusterFile() {
	for {
		time.Sleep(time.Millisecond * 1500)
		nodes := delRepeatElem(clu.Node)
		newClu := Cluster{Node: nodes}
		// 全局变量的更新
		LocalNode = &newClu
		// 写入文件
		SaveClusterFile()
		//fmt.Println(len(newClu.Node))
	}

}
func delRepeatElem(nodes []*Node) []*Node {
	for i := 1; i < len(nodes)-1; i++ {
		if (nodes[i].IP == nodes[i+1].IP && nodes[i].Port == nodes[i+1].Port) ||
			(nodes[i].IP == nodes[i-1].IP && nodes[i].Port == nodes[i+1].Port) {
			nodes = append(nodes[:i], nodes[i+1:]...) //删除重复元素
			i--
		}
	}
	return nodes
}
