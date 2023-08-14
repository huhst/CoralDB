package webBackstage

import (
	"CoralDB/Cluster"
	"CoralDB/GRPC"
	"CoralDB/Raft"
	"CoralDB/consensus"
	"CoralDB/blockchain/blockchain_data"
	"CoralDB/serverExec/service"
	"CoralDB/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var power = 0
var conn *grpc.ClientConn
var client service.ServerClient
var uid string
var userName string
var passWord string

// Admin 主页
func Admin(c *gin.Context) {
	// 获取用户数量，判断用户文件数量
	path := "./userManage/user_files"
	util.DirMap = make(map[string]float64)
	util.Wait.Add(1)
	go util.ScanDir(path, &util.SyncM)
	util.Wait.Wait()
	var fileCount int //文件数量
	util.SyncM.Range(func(key, value interface{}) bool {
		fileCount++
		return true
	})
	// 得到节点IP
	ip := util.LocalIP
	// 得到服务系节点个数
	var nodeSum int
	if Cluster.LocalNode == nil {
		nodeSum = 1
	} else {
		nodeSum = len(Cluster.LocalNode.Node)
	}
	// 得到最新区块id
	id := blockchain_data.LocalDataBlockChain.LastID
	// 得到最新区块的hash
	hash := blockchain_data.LocalDataBlockChain.TailHash
	newBlockHash := fmt.Sprintf("%x", hash)
	// 页面的渲染
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"userSum":      fileCount,
		"serverIp":     ip,
		"serverSum":    nodeSum,
		"newBlockId":   id,
		"newBlockHash": newBlockHash,
	})
}

// IndexPost post请求
func AdminPost(c *gin.Context) {

	json := make(map[string]interface{}) //注意该结构接受的内容
	err := c.BindJSON(&json)
	if err != nil {
		log.Printf(err.Error())
	}
	// 是添加用户
	if util.Strval(json["key"]) == "adduser" {
		userName := util.Strval(json["name"])
		passWord := util.Strval(json["password"])
		// 注册成功
		if Server.RegisterAccount(userName, passWord) {
			messageMap := map[string]interface{}{
				"info": "添加成功",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			messageMap := map[string]interface{}{
				"info": "添加失败，用户已存在，不能创建同名用户",
			}
			c.JSON(http.StatusOK, messageMap)
		}
	}
	// 是创建集群
	if util.Strval(json["key"]) == "creatCluster" {
		// 判断集群文件是否存在
		if util.FileExist("./ClusterInfo") {
			// 如果存在
			messageMap := map[string]interface{}{
				"info": "已加入集群，无法再创建集群",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			alg := consensus.NewAlgorand()
			Cluster.Init(util.Strval(json["password"]))
			// 写入文件
			Cluster.SaveClusterFile()
			// 创建监听
			go Cluster.Server()

			go Raft.Start(&Server.TxPool)
			go alg.Start()

			messageMap := map[string]interface{}{
				"info": "集群创建成功",
			}
			c.JSON(http.StatusOK, messageMap)
		}
	}
	// 是加入集群
	if util.Strval(json["key"]) == "joinCluster" {
		if util.FileExist("./ClusterInfo") {
			// 如果存在
			messageMap := map[string]interface{}{
				"info": "已加入集群，无法再加入集群",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			alg := consensus.NewAlgorand()
			ip := util.Strval(json["ip"])
			port, err := strconv.Atoi("3301")
			if err != nil {
				log.Panic(err)
			}
			key := util.Strval(json["password"])
			err = GRPC.JoinToCluster(ip, port, key)
			if err != nil {
				log.Panic(err)
			}
			// 建立监听
			go Cluster.Server()
			// 广播自己已加入集群
			GRPC.BroadCast(ip, port)
			go Raft.Start(&Server.TxPool)
			go alg.Start()
			// 同步区块链
			time.Sleep(time.Second)
			GRPC.TableBlockSynchronization()
			GRPC.DataBlockSynchronization()
			messageMap := map[string]interface{}{
				"info": "加入集群成功",
			}
			c.JSON(http.StatusOK, messageMap)
		}
	}
	// 设置打包数
	if util.Strval(json["key"]) == "set" {
		numS := util.Strval(json["num"])
		num, err := strconv.Atoi(numS)
		if err != nil {
			messageMap := map[string]interface{}{
				"info": "设置失败",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			Server.TxPool.SetPackNumber(num)
			messageMap := map[string]interface{}{
				"info": "设置成功",
			}
			c.JSON(http.StatusOK, messageMap)
		}

	}

}

// User 主页
func User(c *gin.Context) {
	// 获取信息
	if power == 1 {
		// 已连接服务器
		var nodeSum int
		if Cluster.LocalNode == nil {
			nodeSum = 0
		} else {
			nodeSum = len(Cluster.LocalNode.Node)
		}
		// 得到最新区块id
		id := blockchain_data.LocalDataBlockChain.LastID
		// 得到最新区块的hash
		hash := blockchain_data.LocalDataBlockChain.TailHash
		newBlockHash := fmt.Sprintf("%x", hash)
		c.HTML(http.StatusOK, "user.html", gin.H{
			"UserName":     "",
			"Power":        "1",
			"userAdd":      "用户未登录",
			"serverIp":     util.LocalIP,
			"serverSum":    nodeSum,
			"newBlockId":   id,
			"newBlockHash": newBlockHash,
		})

	}
	if power == 2 {
		var nodeSum int
		if Cluster.LocalNode == nil {
			nodeSum = 0
		} else {
			nodeSum = len(Cluster.LocalNode.Node)
		}
		// 得到最新区块id
		id := blockchain_data.LocalDataBlockChain.LastID
		// 得到最新区块的hash
		hash := blockchain_data.LocalDataBlockChain.TailHash
		newBlockHash := fmt.Sprintf("%x", hash)
		userAdd := Server.CheckAccount(uid)
		c.HTML(http.StatusOK, "user.html", gin.H{
			"UserName":     userName,
			"Power":        "2",
			"userAdd":      userAdd,
			"serverIp":     util.LocalIP,
			"serverSum":    nodeSum,
			"newBlockId":   id,
			"newBlockHash": newBlockHash,
		})

	}
	if power == 0 {
		c.HTML(http.StatusOK, "user.html", gin.H{
			"UserName":     "",
			"Power":        "0",
			"userAdd":      "未连接服务器",
			"serverIp":     "未连接服务器",
			"serverSum":    "未连接服务器",
			"newBlockId":   "未连接服务器",
			"newBlockHash": "未连接服务器",
		})
	}
}

// UserPost post请求
func UserPost(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	err := c.BindJSON(&json)
	if err != nil {
		log.Printf(err.Error())
	}
	// 连接服务器
	if util.Strval(json["key"]) == "connectServer" {
		ip := util.Strval(json["ip"])
		password := util.Strval(json["password"])
		if password != util.ConnPas {
			messageMap := map[string]interface{}{
				"info": "连接密码错误",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			conn, err = grpc.Dial(fmt.Sprintf("%s:%d", ip, 8888), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				messageMap := map[string]interface{}{
					"info": "网络问题，连接失败",
				}
				c.JSON(http.StatusOK, messageMap)
			} else {
				power = 1
				messageMap := map[string]interface{}{
					"info": "连接成功",
				}
				c.JSON(http.StatusOK, messageMap)
			}
			client = service.NewServerClient(conn)
		}
	}

	// 用户登录
	if util.Strval(json["key"]) == "login" {
		userName = util.Strval(json["username"])
		passWord = util.Strval(json["password"])
		// 登录
		login, err := client.Cmd(context.Background(), &service.CommandRequest{LoginAccount: userName, LoginPassword: passWord})
		if err != nil {
			log.Println(err.Error())
		}
		if login.Login == "False" {
			messageMap := map[string]interface{}{
				"info": "登录失败，请检查用户名和密码是否正确",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			power = 2
			uid = userName + "-QAQ-" + passWord
			messageMap := map[string]interface{}{
				"info": "登录成功",
			}
			c.JSON(http.StatusOK, messageMap)
		}
	}
	// 用户退出
	if util.Strval(json["key"]) == "exit" {
		_, err = client.Cmd(context.Background(), &service.CommandRequest{LogoutIng: "true", Uid: uid})
		if err != nil {
			messageMap := map[string]interface{}{
				"info": "退出失败",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			power = 1
			messageMap := map[string]interface{}{
				"info": "退出成功",
			}
			c.JSON(http.StatusOK, messageMap)
		}
	}
	// 添加共享表
	if util.Strval(json["key"]) == "creatTable" {
		tableName := util.Strval(json["tablename"])
		userPower := util.Strval(json["userpower"])
		info, err := client.Cmd(context.Background(), &service.CommandRequest{TabelName: tableName, PermissionTable: userPower, Uid: uid})
		if err != nil {
			log.Printf(err.Error())
		}
		if info.TabelOk == "True" {
			messageMap := map[string]interface{}{
				"info": "创建成功",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			messageMap := map[string]interface{}{
				"info": "创建失败，共享表已存在",
			}
			c.JSON(http.StatusOK, messageMap)
		}

	}
	//修改共享表
	if util.Strval(json["key"]) == "modifyTable" {
		tableName := util.Strval(json["tablename"])
		userPower := util.Strval(json["userpower"])
		info, err := client.Cmd(context.Background(), &service.CommandRequest{TabelName: tableName, PermissionTable: userPower, Uid: uid})
		if err != nil {
			log.Printf(err.Error())
		}
		if info.TabelOk == "True" {
			messageMap := map[string]interface{}{
				"info": "修改成功",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			messageMap := map[string]interface{}{
				"info": "修改失败，无操作权限或共享表不存在",
			}
			c.JSON(http.StatusOK, messageMap)
		}

	}
	// 插入数据
	if util.Strval(json["key"]) == "addDates" {
		dataKey := util.Strval(json["datakey"])
		value := util.Strval(json["value"])
		tableName := util.Strval(json["tablename"])
		info, err := client.Cmd(context.Background(), &service.CommandRequest{TabelName: tableName, Key: dataKey, Value: value, Uid: uid})
		if err != nil {
			log.Printf(err.Error())
		}
		if info.PutOk == "True" {
			messageMap := map[string]interface{}{
				"info": "插入数据成功",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			messageMap := map[string]interface{}{
				"info": "插入数据失败",
			}
			c.JSON(http.StatusOK, messageMap)
		}

	}
	// 修改数据
	if util.Strval(json["key"]) == "modifyDates" {
		dataKey := util.Strval(json["datakey"])
		value := util.Strval(json["value"])
		tableName := util.Strval(json["tablename"])
		info, err := client.Cmd(context.Background(), &service.CommandRequest{TabelName: tableName, Key: dataKey, Value: value, Uid: uid})
		if err != nil {
			log.Printf(err.Error())
		}
		if info.PutOk == "True" {
			messageMap := map[string]interface{}{
				"info": "修改数据成功",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			messageMap := map[string]interface{}{
				"info": "修改数据失败",
			}
			c.JSON(http.StatusOK, messageMap)
		}
	}
	// 断开连接
	if util.Strval(json["key"]) == "close" {
		err := conn.Close()
		if err != nil {
			messageMap := map[string]interface{}{
				"info": "断开连接失败",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			power = 0
			messageMap := map[string]interface{}{
				"info": "断开连接成功",
			}
			c.JSON(http.StatusOK, messageMap)
		}
	}
}

// Data 主页
func Data(c *gin.Context) {
	// 获取信息
	if power == 1 {
		c.HTML(http.StatusOK, "data.html", gin.H{
			"UserName": "",
			"Power":    "1",
		})

	}
	if power == 2 {
		c.HTML(http.StatusOK, "data.html", gin.H{
			"UserName": userName,
			"Power":    "2",
		})

	}
	if power == 0 {
		c.HTML(http.StatusOK, "data.html", gin.H{
			"UserName": "",
			"Power":    "0",
		})
	}

}

// DataPost post请求
func DataPost(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	err := c.BindJSON(&json)
	if err != nil {
		log.Printf(err.Error())
	}
	if util.Strval(json["key"]) == "dataQuery" {
		key := util.Strval(json["datakey"])
		tableName := util.Strval(json["tablename"])

		info := Server.Get(uid, key, tableName, true)
		if info != "" {
			// 获取数据
			slice := strings.Split(info, ",")
			messageMap := map[string]interface{}{
				"info":      "查询成功",
				"tablename": slice[0],
				"key":       slice[1],
				"value":     slice[2],
				"user":      slice[3],
				"time":      slice[4],
				"gettime":   slice[5],
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			messageMap := map[string]interface{}{
				"info": "查询失败，请确保key和表名正确",
			}
			c.JSON(http.StatusOK, messageMap)
		}
	}

}

// Table 主页
func Table(c *gin.Context) {
	// 获取信息
	if power == 1 {
		// 已连接服务器
		c.HTML(http.StatusOK, "table.html", gin.H{
			"UserName": "",
			"Power":    "1",
		})

	}
	if power == 2 {
		c.HTML(http.StatusOK, "table.html", gin.H{
			"UserName": userName,
			"Power":    "2",
		})

	}
	if power == 0 {
		c.HTML(http.StatusOK, "table.html", gin.H{
			"UserName": "",
			"Power":    "0",
		})
	}
}

// TablePost post请求
func TablePost(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	err := c.BindJSON(&json)
	if err != nil {
		log.Printf(err.Error())
	}
	if util.Strval(json["key"]) == "getMyTables" {
		infos := Server.MyTable(uid, true)
		if len(infos) == 0 {
			messageMap := map[string]interface{}{
				"info": "查询失败",
			}
			c.JSON(http.StatusOK, messageMap)
		} else {
			messageMap := map[string]interface{}{
				"info":  "查询成功",
				"sum":   len(infos),
				"infos": infos,
			}
			c.JSON(http.StatusOK, messageMap)
		}

	}
}
