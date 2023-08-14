package userManage

import (
	"CoralDB/util"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
	"os"
	"sync"
)

// 每个节点算区块链P2P的一个节点，里面的用户是依赖于这个节点的。不是独立与这个节点的。
// 用户是绝对信任这个节点的，将对应用户的账户文件存储在其依赖的节点是可靠的。

type UserManage struct {
	EffectiveUser map[string]*Account // map[username][account]
	UserFilesPath string
	Mutex         sync.Mutex
}

func (umg *UserManage) Init() {
	umg.Mutex.Lock()
	defer umg.Mutex.Unlock()

	umg.EffectiveUser = make(map[string]*Account)
	umg.UserFilesPath = "userManage/user_files/"
}

type Account struct {
	UserName   string
	Password   string
	PublicKey  []byte
	PrivateKey ecdsa.PrivateKey
	Address    string
}

// UserManage 是服务器的用户管理模块
// 1.用户注册 - ctime.220903-19:53
// 2.用户登录 - ctime.220903-20:31
// 3.用户查看自己信息 - ctime.220904-16:17
// 4.用户退出登录 - ctime.220904-16:38

// Check 查看在server用户是否是在登录的状态。
func (umg *UserManage) Check(UID string) (Account, error) {
	if _, has := umg.EffectiveUser[UID]; has {
		return *umg.EffectiveUser[UID], nil
	}
	return Account{}, errors.New("user not logged in")
}

// Register 1.用户注册 - 创建一个用户文件到本地（用户名，密码， 私钥）
// err1: 当用户已经被注册时，返回错误-ct
// err2: ecdsa.GenerateKey 执行错误
// err3: 保存账户文件到本地发生错误
func (umg *UserManage) Register(username, password string) error {
	umg.Mutex.Lock()
	defer umg.Mutex.Unlock()

	// 0.检查是否已经存在同名的账户文件
	if util.IsExistFile(umg.UserFilesPath + username + ".act") {
		return errors.New("user has been registered")
	}

	a := new(Account)
	a.UserName = username
	a.Password = password
	// 1.生成账户
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	a.PrivateKey = *privateKey
	a.PublicKey = append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	a.Address = calculateAddress(a.PublicKey)
	// 2.保存到本地
	err = umg.SaveAccount2File(a)
	if err != nil {
		return err
	}
	return nil
}

// SignIn 2.用户登录 - 通过本地用户文件校验用户名和密码。成功后，将用户的信息注册到内存中。
// err1: 在server没有这个用户的用户文件。用户没有注册-ct
// err2: 加载用户文件时发生错误
// err3: 密码错误-ct
// err4: 用户已经处于登录状态-ct
func (umg *UserManage) SignIn(username, password string) error {
	umg.Mutex.Lock()
	defer umg.Mutex.Unlock()

	// 0.检查用户是否存在（用户名是不是正确）
	if !util.IsExistFile(umg.UserFilesPath + username + ".act") {
		return errors.New("the user does not exist")
	}

	// 1.验证密码时候正确
	a := new(Account)
	a.UserName = username
	err := umg.LoadAccount2File(a)
	if err != nil {
		return err
	}
	if a.Password != password {
		return errors.New("password err")
	}

	// 2.添加用户的数据到内存中
	if _, has := umg.EffectiveUser[username+"-QAQ-"+password]; has {
		return errors.New("user is already logged in")
	}
	// 保存ID->Address 到内存中。
	umg.EffectiveUser[username+"-QAQ-"+password] = a

	return nil
}

// ViewAccount 3.用户查看自己的用户信息。返回对应UID的账户。
// err1: 用户没有登录,返回一个空账户和err。
func (umg *UserManage) ViewAccount(UID string) (Account, error) {
	umg.Mutex.Lock()
	defer umg.Mutex.Unlock()

	return umg.Check(UID)
}

// SignOut 4.用户退出登录 - 删除映射
// err1: 用户没有登录
func (umg *UserManage) SignOut(ID string) error {
	umg.Mutex.Lock()
	defer umg.Mutex.Unlock()

	// 用户没有登录
	_, err := umg.Check(ID)
	if err != nil {
		return nil
	}
	// 删除登录信息
	delete(umg.EffectiveUser, ID)
	return nil
}

func (umg *UserManage) SaveAccount2File(account *Account) error {

	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	// 编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(account)
	if err != nil {
		return err
	}
	// 写入文件
	content := buffer.Bytes()
	err = os.WriteFile(umg.UserFilesPath+account.UserName+".act", content, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (umg *UserManage) LoadAccount2File(account *Account) error {
	//判断文件是否存在
	if !util.IsExistFile(umg.UserFilesPath + account.UserName + ".act") {
		return errors.New("user is not exit")
	}
	//读取文件
	content, err := os.ReadFile(umg.UserFilesPath + account.UserName + ".act")
	if err != nil {
		return err
	}
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(account)
	if err != nil {
		return err
	}
	return nil
}

// CalculateAddress 初始化钱包里面的地址
// 根据钱包的公钥私钥得到钱包对应的地址
func calculateAddress(publicKey []byte) string {
	// https://www.cnblogs.com/kumata/p/10477369.html
	publicKeyHash := hashPublicKey(publicKey)

	version := 0x00
	payload := append([]byte{byte(version)}, publicKeyHash...)

	checksum := checkSum(payload)
	payload = append(payload, checksum...)

	address := base58.Encode(payload)

	return address
}

func hashPublicKey(publicKey []byte) []byte {

	sha256Hash := sha256.Sum256(publicKey)

	ripeMD160Hash := ripemd160.New()
	_, err := ripeMD160Hash.Write(sha256Hash[:])
	if err != nil {
		log.Panic(err)
	}
	publicHash := ripeMD160Hash.Sum(nil)

	return publicHash
}

func checkSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	checksum := second[0:4]
	return checksum
}
