package main

import (
	"CoralDB/Cluster"
	"CoralDB/GRPC"
	"CoralDB/Raft"
	"CoralDB/blockqueue"
	"CoralDB/client"
	"CoralDB/consensus"
	"CoralDB/server"
	"CoralDB/serverExec"
	"CoralDB/webBackstage"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	blockqueue.LocalDataBlockQueue = blockqueue.Init()
	blockqueue.LocalTableBlockQueue = blockqueue.Init()
	blockqueue.LocalDataTxQueue = blockqueue.Init()
	//blockqueue.AlgToBC = blockqueue.Init()
	//blockqueue.BCToAlg = blockqueue.Init()
	s := new(server.Server)
	serverExec.RPCs = s
	client.Cserver = s
	s.Init()

	// 运行web
	//go webRun(s)
	// 判断集群文件是否存在，如果不存在则直接执行cmd程序
	_, err := os.Stat("./ClusterInfo")
	if os.IsNotExist(err) {
		//go s.Command()
		go GRPC.DataBlockDistribute()
		go GRPC.TableBlockDistribute()
		go GRPC.SubmitDataTransactions()

		go client.StartClient()
		serverExec.ServerStart()

	} else {
		// 启动Raft, 读取集群文件
		cluster, err := Cluster.LoadClusterFile("./ClusterInfo")
		if err != nil {
			log.Panic(err)
		}

		Cluster.LocalNode = cluster
		for _, node := range cluster.Node {
			fmt.Println(node.IP, node.Port)
		}

		go Cluster.Server()
		// 集群中节点的更新
		go cluster.UpdateClusterFile()

		go Raft.Start(&s.TxPool)

		alg := consensus.NewAlgorand()
		//go Raft.Start(&s.TxPool)
		alg.Start()

		//go s.Command()
		go GRPC.DataBlockDistribute()
		go GRPC.TableBlockDistribute()
		go GRPC.SubmitDataTransactions()

		go client.StartClient()
		serverExec.ServerStart()
	}

}

// Run 启动运行web
func webRun(s *server.Server) {
	// 关闭输出
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()
	// 加载页面
	router.LoadHTMLGlob("templates/*")
	// 加载静态文件
	router.Static("static", "./static")
	// 设置默认路由当访问错时创建一个路由handler
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"msg": "not found",
		})
	})
	webBackstage.InitRouter(router, s)
	err := router.Run(":4000")
	if err != nil {
		return
	}

}
