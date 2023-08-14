package main

import (
	"CoralDB/serverExec/service"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
	"time"
)

//
//import (
//	"context"
//	"flag"
//	"fmt"
//	"github.com/yan259128/CoralDB/serverExec/service"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//	"log"
//	"time"
//)
//
//const (
//	defaultName = "world"
//)
//
//var (
//	addr = flag.String("addr", "192.168.1.101:8888", "the address to connect to")
//	//addr = flag.String("addr", "127.0.0.1:8888", "the address to connect to")
//	//name = flag.String("name", defaultName, "Name to greet")
//	// table a 1FYrDeBmc1yLUkTFuJEtjFpRjjhWQdqspE7
//)
//
//func main() {
//	flag.Parse()
//
//	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		log.Fatalf("did not connect: %v", err)
//	}
//	//defer conn.Close()
//	c := service.NewServerClient(conn)
//	//ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
//
//	username, password := "test", "test"
//	//key := "1M9"
//	//value := "1M9"
//	taableName := "test"
//	//pli := "1M9ZdF3UwmTK7rZATPk5BA23m9V6K4Zn6u7"
//
//	a1, err := c.Cmd(context.Background(), &service.CommandRequest{LoginAccount: username, LoginPassword: password})
//	fmt.Println(a1.Login)
//	if err != nil {
//		log.Fatalf("could not sign in: %v", err)
//	}
//
//	uid := username + "-QAQ-" + password
//
//	//_, err = c.Cmd(context.Background(), &service.CommandRequest{Uid: uid, TabelName: taableName, PermissionTable: pli})
//	//fmt.Println(a1.Login)
//	//if err != nil {
//	//	log.Fatalf("could not sign in: %v", err)
//	//}
//
//	count := 0
//
//	start1 := time.Now() // 获取当前时间
//
//	for i := 0; i < 1000; i++ {
//		key := "key" + string(i)
//		//value := "value" + string(i)
//		//_, err := c.Cmd(context.Background(), &service.CommandRequest{Uid: uid, TabelName: taableName, Key: key, Value: value})
//		_, err := c.Cmd(context.Background(), &service.CommandRequest{Uid: uid, Key: key, TabelName: taableName})
//		if err != nil {
//			fmt.Println(err)
//		}
//		count++
//
//		//fmt.Printf("%d\n", count)
//	}
//	//fmt.Printf("%d  ", count)
//
//	elapsed := time.Since(start1)
//
//	////fmt.Println("在表相关链中查找")
//	//fmt.Println("在lru历史队列中查找")
//	////fmt.Println("在lru缓存队列中查找")
//	fmt.Printf("%d 个读取响应完成耗时：%v\n", count, elapsed)
//	fmt.Println("每秒完成响应", float64(count)/(elapsed.Seconds()), "/s")
//}

var nodeNum = 4
var operationCount = 1
var pkgNum = "2000"

var ipFront = "192.169.10."
var ipAfter = 2
var port = 8888

var nodeIp []string

func main() {
	////fmt.Println(0 % 10)
	cloInit()
	testInit()
	//time.Sleep(time.Second)
	//judgeTable()
	//times := "1689933766133538,1689933766138468,1689933766138611,1689933766146333,1689933766133631,1689933766139155,1689933766139602,1689933766146403,1689933766133684,1689933766142181,1689933766142373,1689933766146180,1689933766133749,1689933766137459,1689933766137832,1689933766146482,1689933766133800,1689933766147592,1689933766147742,1689933766147855,1689933766133959,1689933766137795,1689933766138185,1689933766146607,1689933766133358,1689933766139308,1689933766139506,1689933766146282"
	//times1 := strings.Split(times, ",")
	//for i := 0; i < len(times1); i++ {
	//	u, err := strconv.Atoi(times1[i])
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(time.UnixMicro(int64(u)))
	//}
	//fmt.Println(time.UnixMicro(1689933766133538))
}

func cloInit() {
	for i := 0; i < nodeNum; i++ {
		ip := ipFront + strconv.Itoa(ipAfter+i)
		nodeIp = append(nodeIp, fmt.Sprintf("%s:%d", ip, port))
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("%s:%d 网络异常", ip, port)
		}
		client := service.NewServerClient(conn)
		if ip == "192.169.10.2" {
			//创建集群
			flag, _ := client.Cmd(context.Background(), &service.CommandRequest{Creat: "creat", CreatPassword: "123"})
			fmt.Printf("node %d creat:%s\n", i, flag.CreatOk)
			//conn.Close()
		} else {
			//ip := ipFront + strconv.Itoa(ipAfter+i-1)
			flag, _ := client.Cmd(context.Background(), &service.CommandRequest{Join: "join", JoinIp: "192.169.10.2", JoinPassword: "123"})
			fmt.Printf("node %d join:%s\n", i, flag.JoinOk)
			time.Sleep(time.Second * 2)

		}
		conn.Close()
	}

}

func testInit() {
	for i := 0; i < nodeNum; i++ {
		ip := ipFront + strconv.Itoa(ipAfter+i)
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("%s:%d 网络异常", ip, port)
		}
		client := service.NewServerClient(conn)
		if ip == "192.169.10.2" {
			userName := "a"
			passWord := "a"
			uid := userName + "-QAQ-" + passWord
			login, err := client.Cmd(context.Background(), &service.CommandRequest{LoginAccount: userName, LoginPassword: passWord})
			if err != nil {
				log.Fatalf("could not sign in: %v", err)
			}
			fmt.Println("login: ", login.Login)

			creatTable, err := client.Cmd(context.Background(), &service.CommandRequest{TabelName: "a", PermissionTable: "1FYrDeBmc1yLUkTFuJEtjFpRjjhWQdqspE4", Uid: uid})
			fmt.Printf("node %d creatTable: %s\n", i, creatTable.TabelOk)

			//ver, err := client.Cmd(context.Background(), &service.CommandRequest{ViewMyinfo: "view", Uid: uid})
			//fmt.Println(ver.Uaddress)

			_, err = client.Cmd(context.Background(), &service.CommandRequest{Set: "set", PkgNum: pkgNum})
			//conn.Close()

		} else {
			_, err = client.Cmd(context.Background(), &service.CommandRequest{Set: "set", PkgNum: pkgNum})
			//conn.Close()
		}
		conn.Close()
	}
}

func judgeTable() {
	ipFront := "192.169.10."
	ipAfter := 2
	port := 8888
	for i := 0; i < nodeNum; i++ {
		ip := ipFront + strconv.Itoa(ipAfter+i)
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("%s:%d 网络异常", ip, port)
		}
		client := service.NewServerClient(conn)

		userName := "a"
		passWord := "a"
		uid := userName + "-QAQ-" + passWord

		if ip == "192.169.10.2" {
			ver, _ := client.Cmd(context.Background(), &service.CommandRequest{ViewMyinfo: "view", Uid: uid})
			fmt.Printf("user %s have : %s \n", userName, ver.TabelNames)
		} else {
			login, err := client.Cmd(context.Background(), &service.CommandRequest{LoginAccount: userName, LoginPassword: passWord})
			if err != nil {
				log.Fatalf("could not sign in: %v", err)
			}
			fmt.Println("login: ", login.Login)

			ver, err := client.Cmd(context.Background(), &service.CommandRequest{ViewMyinfo: "view", Uid: uid})
			fmt.Printf("user%d %s have : %s \n", i, userName, ver.TabelNames)
		}
		conn.Close()

	}
}

func testTPS(operationNum int) {
	for i := 0; i < operationNum; i++ {
		id := i % nodeNum
		go func(id int) {
			switch id {
			case 0:

			case 1:
			case 2:
			case 3:
			}
		}(id)
	}

}

func Insert(addr, key, value string) {

}

func read(addr, key, table string) {

}
