package main

import (
	"bytes"
	"fmt"
	bc "github.com/OhYee/blockchain/blockchain"
	gb "github.com/OhYee/goutils/bytes"
	"github.com/xtaci/kcp-go"
	"math/rand"
	"net"
)

var Port int

func startServer() net.Listener {
	var listener net.Listener
	var err error = fmt.Errorf("")

	for err != nil {
		Port = (rand.Intn(55535) % 55535) + 10000
		listener, err = kcp.Listen(fmt.Sprintf("127.0.0.1:%d", Port))
	}

	serverLogger.Printf("Server start at %d\n", Port)

	for i := 0; i < 10; i++ {
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					fmt.Printf("Server error %s\n", err)
				} else {
					server(conn)
				}
			}

		}()
	}

	return listener
}

func server(conn net.Conn) {
	defer conn.Close()
	b, err := gb.ReadNBytes(conn, 1)
	if err != nil {
		serverLogger.Println(err)
	}
	switch b[0] {
	case 'r':
		responseRequest(conn)
	case 's':
		responseSend(conn)
	case 'b':
		responseBlock(conn)
	}
}

func responseRequest(conn net.Conn) {
	serverLogger.Printf("an <request> from %s\n", conn.RemoteAddr())
	defer serverLogger.Printf("<request> close with %s\n", conn.RemoteAddr())

	buf := bytes.NewBuffer([]byte{})
	blocks := blockchain.GetBlocks()
	buf.Write(gb.FromInt32(int32(len(blocks))))
	for _, block := range blocks {
		gb.WriteWithLength32(buf, block.ToBytes())
	}
	conn.Write(buf.Bytes())

	_, err := gb.ReadNBytes(conn, 1)
	if err != nil {
		serverLogger.Println(err)
	}
}

func responseSend(conn net.Conn) {
	serverLogger.Printf("a <send> from %s\n", conn.RemoteAddr())
	defer serverLogger.Printf("<send> close with %s\n", conn.RemoteAddr())

	b, err := gb.ReadWithLength32(conn)
	if err != nil {
		serverLogger.Println(err)
		return
	}
	conn.Write([]byte{0})

	textArray, err := NewTextArrayFromBytes(b)
	if err != nil {
		serverLogger.Println(err)
		return
	}

	for _, t := range textArray.data {
		blockchain.ModifyData(t)
	}
}

func responseBlock(conn net.Conn) {
	serverLogger.Printf("a <block> from %s\n", conn.RemoteAddr())
	defer serverLogger.Printf("<block> close with %s\n", conn.RemoteAddr())

	b, err := gb.ReadWithLength32(conn)
	if err != nil {
		serverLogger.Println(err)
		return
	}

	var block bc.Block
	block, err = blockchain.NewBlockFromBytes(b)

	if !blockchain.AddBlock(block) {
		conn.Write([]byte{1})
		var port int32
		port, err = gb.ReadInt32(conn)
		if err != nil {
			serverLogger.Println(err)
			return
		}
		go func() {
			if err := request(int(port)); err != nil {
				serverLogger.Println(err)
			}
		}()
	}
	conn.Write([]byte{0})

}
