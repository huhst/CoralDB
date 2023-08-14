package blockchain_data

import (
	"CoralDB/MerkleTree"
	"CoralDB/common"
	"CoralDB/util"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Round             uint64         // 区块的ID Round
	CurrentBlockHash  []byte         // 当前区块的 HASH
	PreviousBlockHash []byte         // 上一个区块的 HASH
	MerKelRoot        []byte         // MerKelRoot MerKelRoot
	Transactions      []*Transaction // Transactions 区块的所有交易
	TimeStamp         uint64         // TimeStamp 时间戳

	Author common.Address // 区块提议者的地址
	Seed   []byte
	Proof  []byte

	Type      int8   // 区块的类型
	Signature []byte // 区块的签名
}

// NewBlock _NewBlock
func NewBlock() Block {
	return Block{}
}

// InitBlock _InitBlock
func (block *Block) InitBlock(transactions []*Transaction, previousBlockHash []byte, ID uint64) {
	block.Round = ID
	block.PreviousBlockHash = previousBlockHash
	block.Transactions = transactions
	var MerKelRootData [][]byte
	// 类型转换
	for i := 0; i < len(transactions); i++ {
		MerKelRootData = append(MerKelRootData, transactionsToBytes(*transactions[i]))
	}
	block.MerKelRoot = MerkleTree.GetMerkleRoot(MerKelRootData).Hash
	block.TimeStamp = uint64(time.Now().Unix())
	block.SetBlockHash()
	emptyHash := sha256.Sum256([]byte{})
	block.Seed = emptyHash[:]
	block.Author = common.HashToAddr(emptyHash)
}

// InitGenesisBlock  _InitGenesisBlock
func (block *Block) InitGenesisBlock(transactions []*Transaction, previousBlockHash []byte, ID uint64) {
	block.Round = ID
	block.PreviousBlockHash = previousBlockHash
	block.Transactions = transactions
	var MerKelRootData [][]byte
	// 类型转换
	for i := 0; i < len(transactions); i++ {
		MerKelRootData = append(MerKelRootData, transactionsToBytes(*transactions[i]))
	}
	block.MerKelRoot = MerkleTree.GetMerkleRoot(MerKelRootData).Hash
	block.TimeStamp = uint64(1234567800)
	block.SetBlockHash()
}

// TransactionsToBytes 交易的转换为byte
func transactionsToBytes(transaction Transaction) []byte {
	var data []byte
	data = append(data, transaction.TxID...)
	data = append(data, transaction.DataID...)
	data = append(data, transaction.Table...)
	data = append(data, transaction.Key...)
	data = append(data, transaction.Value...)
	data = append(data, transaction.Possessor...)
	data = append(data, util.Int64ToBytes(transaction.TimeStamp)...)
	data = append(data, transaction.PublicKey...)
	data = append(data, transaction.Signature...)
	return data
}

// SetBlockHash 计算区块的 HASH
func (block *Block) SetBlockHash() {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	blockHash := sha256.Sum256(buffer.Bytes())
	block.CurrentBlockHash = blockHash[:]
}

// Serialize 序列化, 将区块转换成字节流
func (block *Block) Serialize() []byte {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

// Deserialize 反序列化, 将字节流换成区块转
func Deserialize(data []byte) Block {

	var block Block
	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	if err != nil {
		fmt.Println(err)
		fmt.Println("解码区块失败")
		log.Panic(err)
	}

	return block
}

// IsGenesisBlock 判断是否为创世区块
func (block *Block) IsGenesisBlock() bool {
	if block.PreviousBlockHash == nil {
		return true
	}
	return false
}
