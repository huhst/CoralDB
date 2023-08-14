package util

// 定义常用的函数与多次复用的函数

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LocalIP LocalPort 本地的ip和端口号（全局）
var LocalIP = GetLocalIp()
var LocalPort = 3301
var LocalIsAccount = false
var IsDone = true
var SubUser int
var Ycsb = make(map[uint64]bool)
var ConnPas = "123"

// 类型转换的函数

// Int64ToBytes 字节转换
func Int64ToBytes(val int64) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, val); err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// Uint64ToBytes 类型转换
func Uint64ToBytes(val uint64) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, val); err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// Int32ToBytes 类型转换
func Int32ToBytes(val int32) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, val); err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// IntToBytes 类型转换
func IntToBytes(val int) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, val); err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// BytesToInt64 类型转换
func BytesToInt64(val []byte) int64 {
	buf := bytes.NewBuffer(val)
	var data int64
	if err := binary.Read(buf, binary.BigEndian, &data); err != nil {
		log.Panic(err)
	}
	return data
}

// BytesToUint64 类型转换
func BytesToUint64(val []byte) uint64 {
	buf := bytes.NewBuffer(val)
	var data uint64
	if err := binary.Read(buf, binary.BigEndian, &data); err != nil {
		log.Panic(err)
	}
	return data
}

// BytesToInt32 类型转换
func BytesToInt32(val []byte) int32 {
	buf := bytes.NewBuffer(val)
	var data int32
	if err := binary.Read(buf, binary.BigEndian, &data); err != nil {
		log.Panic(err)
	}
	return data
}

// BytesToInt 类型转换
func BytesToInt(val []byte) int {
	buf := bytes.NewBuffer(val)
	var data int
	if err := binary.Read(buf, binary.BigEndian, &data); err != nil {
		log.Panic(err)
	}
	return data
}

// StrToBytes 类型转换
func StrToBytes(val string) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, val); err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// BytesToStr 类型转换
func BytesToStr(val []byte) string {
	return string(val)
}

// GetLocalIp 获取本机的ip
func GetLocalIp() (ip string) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Panic(err)
		}
	}(conn) // 链接的关闭
	if err != nil {
		log.Panic(err)
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)   // 获取回环地址，既可以与外界通信的ip地址
	ip = strings.Split(localAddr.String(), ":")[0] // 得到ip
	//port = strings.Split(localAddr.String(), ":")[1]
	return ip

	//公网ip
	//resp, err := http.Get("http://myexternalip.com/raw")
	//if err != nil {
	//	return ""
	//}
	//defer resp.Body.Close()
	//content, _ := ioutil.ReadAll(resp.Body)
	//return string(content)
}
func GetLocalIp1() (ip string) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Panic(err)
		}
	}(conn) // 链接的关闭
	if err != nil {
		log.Panic(err)
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)   // 获取回环地址，既可以与外界通信的ip地址
	ip = strings.Split(localAddr.String(), ":")[0] // 得到ip
	//port = strings.Split(localAddr.String(), ":")[1]
	return ip
}

// AesCTREncrypt 加密
func AesCTREncrypt(src, key []byte) []byte {
	//fmt.Printf("明文： %s\n", src)
	//1. 创建一个cipher.Block接口。
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	//fmt.Println("aes block size : ", block.BlockSize())

	iv := bytes.Repeat([]byte("1"), block.BlockSize())

	//2. 选择分组模式：ctr
	stream := cipher.NewCTR(block, iv)

	//3. 加密操作
	stream.XORKeyStream(src /*密文*/, src /*明文*/)

	return src
}

// AesCTRDecrypt 解密
func AesCTRDecrypt(cipherData, key []byte) []byte {

	//1. 创建一个cipher.Block接口。
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	iv := bytes.Repeat([]byte("1"), block.BlockSize())

	//2. 选择分组模式：ctr
	stream := cipher.NewCTR(block, iv)

	//3. 解密操作
	stream.XORKeyStream(cipherData /*明文*/, cipherData)

	return cipherData
}

// IsExistFile file 是不是 存在
func IsExistFile(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

// ConvertStrSlice2Map 将字符串 slice 转为 map[string]struct{}。
func ConvertStrSlice2Map(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// InMap 判断字符串是否在 map 中。
func InMap(m map[string]struct{}, s string) bool {
	_, ok := m[s]
	return ok
}

func DialCustom(network, address string, timeout time.Duration, localIP []byte, localPort int) (net.Conn, error) {
	netAddr := &net.TCPAddr{Port: localPort}

	if len(localIP) != 0 {
		netAddr.IP = localIP
	}
	d := net.Dialer{Timeout: timeout, LocalAddr: netAddr}
	return d.Dial(network, address)
}

var lock sync.Mutex

func SetYcsb(round uint64, flag bool) {
	lock.Lock()
	defer lock.Unlock()
	round = round % 10
	Ycsb[round] = flag
}
func GetYcsb(round uint64) bool {
	lock.Lock()
	defer lock.Unlock()
	round = round % 10
	return Ycsb[round]
}

// 获取文件数量
var SyncM sync.Map
var DirMap map[string]float64
var Wait sync.WaitGroup

func ScanDir(path string, syncM *sync.Map) {
	defer Wait.Done()
	dirAry, err := ioutil.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}
	for _, e := range dirAry {
		if e.IsDir() {
			Wait.Add(1)
			go ScanDir(filepath.Join(path, e.Name()), syncM)
		} else {
			syncM.Store(filepath.Join(path, e.Name()), (float64(e.Size()) / 1024))
		}
	}

}

// Strval 类型转换
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

//判断文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
