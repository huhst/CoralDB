//package blockchain_table
//
//import (
//	"bytes"
//	"crypto/ecdsa"
//	"crypto/elliptic"
//	"crypto/rand"
//	"crypto/sha256"
//	"encoding/gob"
//	"log"
//	"math/big"
//	"time"
//)
//
//// Transaction
//// TxID   		交易的ID,HASH(交易的所有信息)
//// Signature 	签名,sign(交易ID)
//type Transaction struct {
//	TxID []byte
//
//	Table string // 唯一的。一个数据库里面不能存在同名的表。
//	// 表的权限信息。
//	// 1.只读权限 ReadOnly 1
//	// 2.读写权限 ReadWrite 2
//	// 3.覆盖写权限 Overwrite 3
//	// 4.表修改权限（更新表权限信息等） TableManger 4
//	PermissionTable []string
//	Possessor       string // 谁发布了这个表
//	TimeStamp       int64  // 交易在本地生成的时间戳. 是在区块链中的生效日期。
//	// 验证信息
//	PublicKey []byte // 交易所有者的公钥
//	Signature []byte // 交易所有者的签名
//}
//
//func (tx *Transaction) Init(table string, permissionTable []string,
//	possessor string, publicKey []byte, privateKey ecdsa.PrivateKey) {
//	// init
//	tx.Table = table
//	tx.PermissionTable = permissionTable
//	tx.Possessor = possessor // 这条数据的所有者。
//	tx.TimeStamp = time.Now().Unix()
//	tx.PublicKey = publicKey
//
//	tx.SetTxID()         // nil
//	tx.Sign(&privateKey) // nil
//}
//
//// SetTxID 交易的标识ID，表明这笔交易的唯一性。
//func (tx *Transaction) SetTxID() {
//	var buffer bytes.Buffer
//	encoder := gob.NewEncoder(&buffer)
//	err := encoder.Encode(transactionsToBytes(*tx))
//	if err != nil {
//		log.Panic(err)
//	}
//	txHash := sha256.Sum256(buffer.Bytes())
//	tx.TxID = txHash[:]
//}
//
//// Sign 由交易发起者提供私钥，为交易签名。
//func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) {
//	signDataHash := tx.TxID
//	r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
//	if err != nil {
//		log.Panic(err)
//	}
//	signature := append(r.Bytes(), s.Bytes()...)
//	tx.Signature = signature
//}
//
//// VerifyTransaction 验证交易是否在传输的时候被篡改。
//// 只需要验证公钥私钥是不是一队就可以
//func VerifyTransaction(txVerify Transaction) bool {
//	// 得到签名, 公钥
//	signature := txVerify.Signature
//	publicKey := txVerify.PublicKey
//	// 1. 验证交易签名的公钥,私钥是不是一对
//	// 签名还原
//	r := big.Int{}
//	s := big.Int{}
//	rData := signature[:len(signature)/2]
//	sData := signature[len(signature)/2:]
//	r.SetBytes(rData)
//	s.SetBytes(sData)
//	// 公钥还原
//	x := big.Int{}
//	y := big.Int{}
//	xData := publicKey[:len(publicKey)/2]
//	yData := publicKey[len(publicKey)/2:]
//	x.SetBytes(xData)
//	y.SetBytes(yData)
//
//	// 得到待验证的hash
//	txVerify.Signature = nil
//	txVerify.TxID = nil
//	txVerify.SetTxID()
//
//	curve := elliptic.P256()
//	ecdsaPublicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
//	if ecdsa.Verify(&ecdsaPublicKey, txVerify.TxID, &r, &s) == false {
//		return false
//	}
//	return true
//}

package blockchain_table

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
	"time"
)

// Transaction
// TxID   		交易的ID,HASH(交易的所有信息)
// Signature 	签名,sign(交易ID)
type Transaction struct {
	TxID []byte

	Table string // 唯一的。一个数据库里面不能存在同名的表。
	// 表的权限信息。
	// 1.只读权限 ReadOnly 1
	// 2.读写权限 ReadWrite 2
	// 3.覆盖写权限 Overwrite 3
	// 4.表修改权限（更新表权限信息等） TableManger 4
	PermissionTable []string
	Possessor       string // 谁发布了这个表
	TimeStamp       int64  // 交易在本地生成的时间戳. 是在区块链中的生效日期。
	// 验证信息
	PublicKey []byte // 交易所有者的公钥
	Signature []byte // 交易所有者的签名
}

func (tx *Transaction) Init(table string, permissionTable []string, possessor string, publicKey []byte, privateKey ecdsa.PrivateKey) {
	// init
	tx.Table = table
	tx.PermissionTable = permissionTable
	tx.Possessor = possessor // 这条数据的所有者。
	tx.TimeStamp = time.Now().Unix()
	tx.PublicKey = publicKey

	tx.SetTxID()         // nil
	tx.Sign(&privateKey) // nil
}

// SetTxID 交易的标识ID，表明这笔交易的唯一性。
func (tx *Transaction) SetTxID() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(transactionsToBytes(*tx))
	if err != nil {
		log.Panic(err)
	}
	txHash := sha256.Sum256(buffer.Bytes())

	tx.TxID = txHash[:]
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) {
	signDataHash := tx.TxID
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, signDataHash[:])
	if err != nil {
		log.Panic(err)
	}
	tx.Signature = sig
}

func VerifyTransaction(txVerify Transaction) bool {
	// 得到签名, 公钥
	signature := txVerify.Signature
	publicKey := txVerify.PublicKey

	// 公钥还原
	x := big.Int{}
	y := big.Int{}
	xData := publicKey[:len(publicKey)/2]
	yData := publicKey[len(publicKey)/2:]
	x.SetBytes(xData)
	y.SetBytes(yData)

	// 得到待验证的hash
	txVerify.Signature = nil
	txVerify.TxID = nil
	txVerify.SetTxID()

	curve := elliptic.P256()
	ecdsaPublicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	if ecdsa.VerifyASN1(&ecdsaPublicKey, txVerify.TxID, signature) == false {
		return false
	}
	return true
}
