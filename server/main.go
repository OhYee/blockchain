package main

import (
	"bufio"
	bc "github.com/OhYee/blockchain/blockchain"
	"os"
	"strconv"
	"strings"

	"github.com/OhYee/rainbow/color"
	"github.com/OhYee/rainbow/log"
)

var (
	blockchain   = bc.NewBlockChain(NewTextArray())
	reader       = bufio.NewReader(os.Stdin)
	serverLogger = log.New().SetColor(
		color.New().SetFrontBlue(),
	).SetOutputToStdout().SetPrefix(
		func() func(string) string {
			serverColor := color.New().SetFrontBlue().SetFontBold()
			return func(s string) string {
				return serverColor.Colorful("[SERVER] ")
			}
		}(),
	)
	systemLogger = log.New().SetColor(
		color.New().SetFrontRed(),
	).SetOutputToStdout().SetPrefix(
		func() func(string) string {
			systemColor := color.New().SetFrontRed().SetFontBold()
			return func(s string) string {
				return systemColor.Colorful("[SYSTEM] ")
			}
		}(),
	)
	uiLogger = log.New().SetOutputToStdout().SetPrefix(func() func(string) string {
		systemColor := color.New().SetFontBold()
		return func(s string) string {
			return systemColor.Colorful("[UI] ")
		}
	}(),
	)
)

func ui() {
	uiLogger.Printf("=================================================================\n")
	uiLogger.Printf("key                          generate a pair of key (pu pr)      \n")
	uiLogger.Printf("add <pu> <pr> <data>         add the <data> to unconfirmed data  \n")
	uiLogger.Printf("generate                     generate a new block                \n")
	uiLogger.Printf("unconfirm                    show all of the unconfirmed data    \n")
	uiLogger.Printf("print                        show blockchain                     \n")
	uiLogger.Printf("hash <hash>                  show block which hash is <hash>     \n")
	uiLogger.Printf("index <index>                show block which index is <index>   \n")
	uiLogger.Printf("clean                        clean the unconfirmed data          \n")
	uiLogger.Printf("exit                         exit the propram                    \n")
	uiLogger.Printf("=================================================================\n")
	uiLogger.Printf("send <port>                  send the unconfirmed data to <port> \n")
	uiLogger.Printf("request <port>               request the blockchain in <port>    \n")
	uiLogger.Printf("block <port>                 send the valid block to <port>      \n")
	uiLogger.Printf("=================================================================\n")

	for {
		uiLogger.Printf("Input your command:\n")
		var command string
		command, err := reader.ReadString('\n')
		if err != nil {
			uiLogger.Println(err)
			continue
		}
		commands := strings.Split(command[:len(command)-1], " ")
		if len(commands) >= 1 {
			switch commands[0] {
			case "key", "k":
				publicKey, privateKey, err := generateKeys()
				if err != nil {
					uiLogger.Printf("Error %s\n", err)
				} else {
					uiLogger.Printf("\npublic key:  %s\nprivate key: %s\n", publicKey.String(), privateKey.String())
				}
			case "add", "a":
				if len(commands) < 3 {
					uiLogger.Printf("using 'add <pu> <pr> <data>' or 'a <pu> <pr> <data>'\n")
				} else {
					add(commands[1], commands[2], strings.Join(commands[3:], " "))
				}
			case "hash", "h":
				if len(commands) == 1 {
					uiLogger.Printf("using 'hash <hash>' or 'h <hash>'\n")
				} else {
					block, ok := getByHash(commands[1])
					if ok {
						uiLogger.Println(block.String())
					} else {
						uiLogger.Printf("Can not find this block\n")
					}
				}
			case "index", "i":
				if len(commands) == 1 {
					uiLogger.Printf("using 'index <index>' or 'i <index>'\n")
				} else {
					n32, err := strconv.ParseInt(commands[1], 10, 32)
					if err != nil {
						uiLogger.Printf("index %s is not a number\n", commands[1])
					} else {
						block, ok := getByIndex(int(n32))
						if ok {
							uiLogger.Println(block.String())
						} else {
							uiLogger.Printf("Can not find this block\n")
						}
					}
				}
			case "clean", "c":
				clean()
			case "generate", "g":
				generate()
			case "unconfirm", "u":
				uiLogger.Println(blockchain.UnconfirmedData.String(""))
			case "print", "p":
				length := blockchain.GetLength()
				for idx, block := range blockchain.GetBlocks() {
					uiLogger.Printf("%d/%d %s", idx+1, length, block.String())
				}
			case "request", "r":
				if len(commands) == 1 {
					uiLogger.Printf("using 'request <port>' or 'r <port>'\n")
				} else {
					n, err := strconv.ParseInt(commands[1], 10, 32)
					if err != nil {
						uiLogger.Println(err)
					} else {
						if err := request(int(n)); err != nil {
							systemLogger.Println(err)
						}
					}
				}
			case "send", "s":
				if len(commands) == 1 {
					uiLogger.Printf("using 'send <port>' or 's <port>'\n")
				} else {
					n, err := strconv.ParseInt(commands[1], 10, 32)
					if err != nil {
						uiLogger.Println(err)
					} else {
						if err := send(int(n)); err != nil {
							systemLogger.Println(err)
						}
					}
				}
			case "exit", "e":
				os.Exit(0)
			case "block", "b":
				if len(commands) == 1 {
					uiLogger.Printf("using 'send <port>' or 's <port>'\n")
				} else {
					n, err := strconv.ParseInt(commands[1], 10, 32)
					if err != nil {
						uiLogger.Println(err)
					} else {
						if err := block(int(n)); err != nil {
							systemLogger.Println(err)
						}
					}
				}
			default:
				uiLogger.Printf("Can not understand '%s' and what you want to do\n", command)
			}
		}
	}
}

func main() {
	listener := startServer()
	defer listener.Close()
	ui()
}
