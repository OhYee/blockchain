package main

import (
	"bufio"
	"bytes"
	"fmt"
	bc "github.com/OhYee/blockchain/blockchain"
	"os"
	"strconv"
	"strings"
	"time"
)

type textArray struct {
	data []string
}

func newStringSlice() *textArray {
	return &textArray{
		data: make([]string, 0),
	}
}

func (text *textArray) Copy() bc.BlockData {
	s := newStringSlice()
	s.data = make([]string, len(text.data))
	copy(s.data, text.data)
	return s
}

func (text *textArray) Reset() {
	text.data = text.data[:0]
}

func (text *textArray) Modify(args ...interface{}) {
	for _, v := range args {
		s, ok := v.(string)
		if ok {
			text.data = append(text.data, s)
		}
	}
}

func (text *textArray) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	for _, s := range text.data {
		buf.WriteString(s)
	}
	return buf.Bytes()
}

func main() {
	var blockchain = bc.NewBlockChain(newStringSlice())
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("=====================================================\n")
	fmt.Printf("add <data>         add the <data> to unconfirm data  \n")
	fmt.Printf("confirm            confirm all of the data (mine)    \n")
	fmt.Printf("unconfirm          show all of the unconfirm data    \n")
	fmt.Printf("show               show blockchain                   \n")
	fmt.Printf("hash <hash>        show block which hash is <hash>   \n")
	fmt.Printf("index <index>      show block which index is <index> \n")
	fmt.Printf("exit               exit the propram                  \n")
	fmt.Printf("=====================================================\n")

	for {
		fmt.Printf("\n\nInput your command:\n")
		var command string
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		commands := strings.Split(command[:len(command)-1], " ")
		if len(commands) >= 1 {
			switch commands[0] {
			case "add", "a":
				if len(commands) == 1 {
					fmt.Printf("using 'add <data>' or 'a <data>'\n")
				} else {
					blockchain.ModifyData(strings.Join(commands[1:], " "))
				}
			case "hash", "h":
				if len(commands) == 1 {
					fmt.Printf("using 'hash <hash>' or 'h <hash>'\n")
				} else {
					block, ok := blockchain.GetBlockOfHash(commands[1])
					if ok {
						fmt.Printf("block %s\n", block.Hash().String())
						fmt.Printf("\tTimestamp: %s\n", time.Unix(block.GetTimestamp(), 0).Format("2006-01-02 15:04:05"))
						fmt.Printf("\tProof: %d\n", block.GetProof())
						fmt.Printf("\tPreHash: %s\n", block.GetPreHash().String())
						ss := block.GetBlockData().(*textArray)
						for idx2, s := range ss.data {
							fmt.Printf("\t%d %s\n", idx2+1, s)
						}
						fmt.Printf("\n")
					} else {
						fmt.Printf("Can not find this block\n")
					}
				}
			case "index", "i":
				if len(commands) == 1 {
					fmt.Printf("using 'index <index>' or 'i <index>'\n")
				} else {
					n32, err := strconv.ParseInt(commands[1], 10, 32)
					if err != nil {
						fmt.Printf("index %s is not a number\n", commands[1])
					} else {
						block, ok := blockchain.GetBlockOfIndex(int(n32))
						if ok {
							fmt.Printf("block %s\n", block.Hash().String())
							fmt.Printf("\tTimestamp: %s\n", time.Unix(block.GetTimestamp(), 0).Format("2006-01-02 15:04:05"))
							fmt.Printf("\tProof: %d\n", block.GetProof())
							fmt.Printf("\tPreHash: %s\n", block.GetPreHash().String())
							ss := block.GetBlockData().(*textArray)
							for idx2, s := range ss.data {
								fmt.Printf("\t%d %s\n", idx2+1, s)
							}
							fmt.Printf("\n")
						} else {
							fmt.Printf("Can not find this block\n")
						}
					}
				}
			case "confirm", "c":
				block := blockchain.NewBlock()
				block = block.Mine()
				blockchain.AddBlock(block)
			case "unconfirm", "u":
				ss := (blockchain.UnconfirmedData).(*textArray)
				for _, s := range ss.data {
					fmt.Printf("\t%s", s)
				}
			case "show", "s":
				length := blockchain.GetLength()
				for idx, block := range blockchain.GetBlocks() {
					fmt.Printf("%d/%d block %s\n", idx+1, length, block.Hash().String())
					fmt.Printf("\tTimestamp: %s\n", time.Unix(block.GetTimestamp(), 0).Format("2006-01-02 15:04:05"))
					fmt.Printf("\tProof: %d\n", block.GetProof())
					fmt.Printf("\tPreHash: %s\n", block.GetPreHash().String())
					ss := block.GetBlockData().(*textArray)
					for idx2, s := range ss.data {
						fmt.Printf("\t%d %s\n", idx2+1, s)
					}
					fmt.Printf("\n")
				}
			case "exit", "e":
				os.Exit(0)
			default:
				fmt.Printf("Can not understand '%s' and what you want to do\n", command)
			}
		}

	}
}
