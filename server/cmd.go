package server

import (
	"CoralDB/Cluster"
	"CoralDB/GRPC"
	"CoralDB/Raft"
	"CoralDB/consensus"
	"CoralDB/util"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const Usage0 = `Help info
  creatcluster joinkey -- 创建集群与设置加入集群的密钥
  joincluster ip port joinkey -- 向已在集群中的节点发出加入集群的请求
  register username userPassword  -- 用户注册
  login username userPassword -- 用户登录
  address -- 查看用户的地址
  table tableName userAddress+permissions -- 创建新的共享表
  put key value tableName -- 向共享表中添加数据
  get key tableName -- 在共享表中查询数据
  gethistory key tableName -- 查询表的更新历史
  mytables -- 查看自己所在的共享表
  root-tables -- 查看自己所在表的权限
  isaccount -- 查看节点是否拥有记账权
  setpkg num -- 设置打包数量
  uin username userpaaword -- 在终端登录用户
  exit -- 退出登录或退出程序
  help -- 输出辅助信息
`

func (s *Server) Command() {

	username, password, uID := "", "", ""

	reader := bufio.NewReader(os.Stdin)

	// 输出自己的IP
	fmt.Printf("监听地址为%s:%d\n", util.LocalIP, util.LocalPort)
	alg := consensus.NewAlgorand()

	for {
		go GRPC.DataBlockDistribute()
		go GRPC.TableBlockDistribute()
		if username != "" {
			fmt.Printf(" %s ", username)
		}
		fmt.Printf("> ")
		args := make([]string, 0)
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Panic(err)
		}
		split := bytes.Split(line, []byte(" "))
		for _, argBytes := range split {
			if string(argBytes) != "" {
				args = append(args, string(argBytes))
			}
		}
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		// 设置在终端显示的用户。 设置在终端显示的用户，并设置对应的用户名，密码，UID。
		case "uin":
			if len(args) == 3 {
				username, password, uID = args[1], args[2], args[1]+"-QAQ-"+args[2]
			} else {
				fmt.Println("user username password")
			}
		// 注册用户到 server
		case "register":
			if len(args) == 3 {
				s.RegisterAccount(args[1], args[2])
			} else {
				fmt.Println("register username password")
			}
		// 登录用户到 server
		case "login":
			if len(args) == 3 {
				s.SignInAccount(args[1], args[2])
			} else {
				fmt.Println("login username password")
			}
		// 查看用户的地址
		case "address":
			if len(args) == 1 {
				s.CheckAccount(uID)
			}
		// u_in , 退出登录。 或者退出程序。
		case "exit":
			if username == "" {
				os.Exit(1)
			}
			if len(args) == 1 {
				s.SignOutAccount(username + "-QAQ-" + password)
				username = ""
				password = ""
			}
			//
			//
		// 创建一个表
		case "table":
			if len(args) >= 3 {
				plist := make([]string, 0)
				for i := 2; i < len(args); i++ {
					plist = append(plist, args[i])
				}
				s.Table(username+"-QAQ-"+password, args[1], plist)
			} else {
				fmt.Println("table tablename 用户地址+用户权限（1-4）...")
			}
		case "put":
			if len(args) == 4 {
				s.Put(username+"-QAQ-"+password, args[1], args[2], args[3])
			}
		case "get":
			if len(args) == 3 {
				s.Get(username+"-QAQ-"+password, args[1], args[2], false)
			}
		case "gethistory":
			if len(args) == 3 {
				s.GetHistory(username+"-QAQ-"+password, args[1], args[2])
			}
		case "gettable":
			if len(args) == 2 {
				s.GetTableData(username+"-QAQ-"+password, args[1])
			}
		case "mytables":
			s.MyTable(uID, false)

		case "creat":
			if len(args) == 2 {
				if util.FileExist("./ClusterInfo") {
					// 如果存在
					fmt.Println("已加入集群，无法再创建集群")
				} else {
					Cluster.Init(args[1])
					// 写入文件
					Cluster.SaveClusterFile()
					// 创建监听
					go Cluster.Server()

					go Raft.Start(&s.TxPool)
					go alg.Start()
				}
			}
		case "join":
			if len(args) == 4 {
				_, err := os.Stat("./ClusterInfo")
				if os.IsNotExist(err) {
					ip := args[1]
					port, err := strconv.Atoi(args[2])
					if err != nil {
						log.Panic(err)
					}
					key := args[3]
					err = GRPC.JoinToCluster(ip, port, key)
					if err != nil {
						log.Panic(err)
					}
					// 建立监听
					go Cluster.Server()
					// 广播自己已加入集群
					GRPC.BroadCast(ip, port)
					go Raft.Start(&s.TxPool)
					go alg.Start()
					// 同步区块链
					time.Sleep(time.Second)
					GRPC.TableBlockSynchronization()
					GRPC.DataBlockSynchronization()

				} else {
					fmt.Println("节点已在集群中")
				}
			}
		case "isaccount":
			if len(args) == 1 {
				fmt.Println(util.LocalIsAccount)
			}
		case "help":
			if len(args) == 1 {
				fmt.Printf(Usage0)
			}
		case "setpkg":
			num, _ := strconv.Atoi(args[1])
			s.TxPool.SetPackNumber(num)

		//	--------------------------------------------------------------------------------------------------------------------
		case "root-tables":
			s.Cache.ROOT_PooledTables()
		case "root-lru":
			s.Cache.ROOT_lruQ()
		case "id":
			if len(args) == 1 {
				fmt.Printf(Usage0)
			}
			a, _ := strconv.Atoi(args[1])
			b := s.dataChain.GetByRound(uint64(a))
			fmt.Println(b.Round)
		case "gen":
			fmt.Println(s.dataChain.GetBlockByHash([]byte{}))

		//case "root-tlist": // 查看表的信息
		//	s.root_tlist()
		//case "root-cache": // 查看缓存队列
		//	s.root_cache()
		//case "root-mod":
		//	s.TxPool.SetMod(util.LocalIsAccount)
		case "root-mod":
			s.TxPool.SetMod1()
		case "printcluster":
			if len(args) == 1 {
				_, err := os.Stat("./ClusterInfo")
				if os.IsNotExist(err) {
					fmt.Println("集群文件不存在")
				} else {
					for _, node := range Cluster.LocalNode.Node {
						fmt.Println(node.IP + ":" + strconv.Itoa(node.Port))
					}
				}
			}

		default:
			fmt.Printf(Usage0)
		}
	}
}
