package consensus

import (
	"CoralDB/blockchain/blockchain_data"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

// Serialize 序列化, 将区块转换成字节流
func Serialize(block *blockchain_data.Block) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

// Deserialize 反序列化, 将字节流换成区块转
func Deserialize(data []byte) blockchain_data.Block {

	var block blockchain_data.Block
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
