package main

import (
	"fmt"
	bc "github.com/OhYee/blockchain/blockchain"
	"github.com/OhYee/cryptography_and_network_security/RSA"
	gb "github.com/OhYee/goutils/bytes"
	"github.com/xtaci/kcp-go"
	"strings"
)

func add(publicKey string, privateKey string, data string) {
	pu := bc.NewHashCodeFromString(strings.Trim(publicKey, " "))
	pr := bc.NewHashCodeFromString(strings.Trim(privateKey, " "))
	text := NewText(pu, pr, data)
	blockchain.ModifyData(text)
}

func getByHash(hash string) (block *bc.Block, ok bool) {
	block, ok = blockchain.GetBlockOfHash(hash)
	return
}

func getByIndex(index int) (block *bc.Block, ok bool) {
	block, ok = blockchain.GetBlockOfIndex(index)
	return
}

func generate() {
	block := blockchain.NewBlock()
	block = block.Mine()
	if blockchain.AddBlock(block) {
		systemLogger.Println("A new block generated")
	} else {
		systemLogger.Println("block invalid")
	}
}

func clean() {
	blockchain.UnconfirmedData.Reset()
	systemLogger.Println("Unconfirmed data cleaned")
}

func request(port int) (err error) {
	conn, err := kcp.Dial(fmt.Sprintf("127.0.0.1:%d", port))
	defer conn.Close()
	if err != nil {
		return
	}
	systemLogger.Printf("<request> command begin with %s\n", conn.RemoteAddr())
	defer systemLogger.Printf("<request> command end with %s\n", conn.RemoteAddr())

	conn.Write([]byte{'r'})

	length32, err := gb.ReadInt32(conn)
	if err != nil {
		return
	}
	length := int(length32)

	if length <= blockchain.GetLength() {
		return
	}
	tempBlockChain := bc.NewBlockChain(NewTextArray())

	for i := 0; i < length; i++ {
		var b []byte
		b, err = gb.ReadWithLength32(conn)
		if err != nil {
			return
		}

		var block bc.Block
		block, err = blockchain.NewBlockFromBytes(b)
		if err != nil {
			return
		}

		tempBlockChain.AddBlock(block)
	}
	conn.Write([]byte{0})

	blockchain.Update(tempBlockChain)
	systemLogger.Printf("update blockchain from %s\n", conn.RemoteAddr())
	return
}

func generateKeys() (pu bc.HashCode, pr bc.HashCode, err error) {
	privateKey, publicKey, err := rsa.Generate()
	if err != nil {
		return
	}
	systemLogger.Printf("Key pair(%v,%v)\n", privateKey, publicKey)
	
	pu = bc.NewHashCodeFromBytes(publicKey)
	pr = bc.NewHashCodeFromBytes(privateKey)
	return
}

func send(port int) (err error) {
	conn, err := kcp.Dial(fmt.Sprintf("127.0.0.1:%d", port))
	defer conn.Close()
	if err != nil {
		return
	}
	systemLogger.Printf("<send> command begin with %s\n", conn.RemoteAddr())
	defer systemLogger.Printf("<send> command end with %s\n", conn.RemoteAddr())

	conn.Write([]byte{'s'})
	gb.WriteWithLength32(conn, blockchain.UnconfirmedData.ToBytes())

	_, err = gb.ReadNBytes(conn, 1)
	if err != nil {
		return
	}

	systemLogger.Printf("Unconfirm data was send to %s\n", conn.RemoteAddr())
	return
}

func block(port int) (err error) {
	conn, err := kcp.Dial(fmt.Sprintf("127.0.0.1:%d", port))
	defer conn.Close()
	if err != nil {
		return
	}
	systemLogger.Printf("<block> command begin with %s\n", conn.RemoteAddr())
	systemLogger.Printf("<block> command end with %s\n", conn.RemoteAddr())

	conn.Write([]byte{'b'})
	gb.WriteWithLength32(conn, blockchain.GetBlocks()[blockchain.GetLength()-1].ToBytes())
	systemLogger.Printf("Send the last block to %s\n", conn.RemoteAddr())

	var b []byte
	b, err = gb.ReadNBytes(conn, 1)
	if err != nil {
		return
	}
	if b[0] != 0 {
		conn.Write(gb.FromInt32(int32(Port)))
		systemLogger.Printf("Send the local port to %s\n", conn.RemoteAddr())
		_, err = gb.ReadNBytes(conn, 1)
		if err != nil {
			return
		}
	}

	return
}
