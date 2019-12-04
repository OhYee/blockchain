package blockchain

import (
	"bytes"
	gb "github.com/OhYee/goutils/bytes"
	"sync"
	"time"
)

type BlockChain struct {
	UnconfirmedData BlockData         // the data wait for add to the blockchain
	blocks          []Block           // blocks of the blockchain
	hashTable       map[string]*Block // a map for hash to block
	mutex           *sync.Mutex
}

func NewBlockChain(initData BlockData) *BlockChain {
	initBlock := NewBlock(1575449947, initData, 0, NewHashCodeFromBytes([]byte{}))
	initBlock = initBlock.Mine()

	initData.Reset()
	return &BlockChain{
		UnconfirmedData: initData,
		blocks:          []Block{initBlock},
		hashTable:       map[string]*Block{initBlock.Hash().String(): &initBlock},
		mutex:           new(sync.Mutex),
	}
}

// ModifyData modify the data
func (bc *BlockChain) ModifyData(args ...interface{}) {
	bc.Lock()
	defer bc.Unlock()
	bc.UnconfirmedData.Modify(args...)
}

// NewBlock init the block
func (bc *BlockChain) NewBlock() Block {
	return NewBlock(time.Now().Unix(), bc.UnconfirmedData, 0, bc.blocks[bc.GetLength()-1].Hash())
}

func (bc *BlockChain) NewBlockFromBytes(b []byte) (block Block, err error) {
	buf := bytes.NewBuffer(b)

	var timestamp int64
	var proof uint64
	var blockData, preHash []byte
	var data BlockData

	if timestamp, err = gb.ReadInt64(buf); err != nil {
		return
	}

	if blockData, err = gb.ReadWithLength32(buf); err != nil {
		return
	}

	if proof, err = gb.ReadUint64(buf); err != nil {
		return
	}

	if preHash, err = gb.ReadWithLength32(buf); err != nil {
		return
	}

	if data, err = bc.UnconfirmedData.FromBytes(blockData); err != nil {
		return
	}

	block = NewBlock(timestamp, data, proof, NewHashCodeFromBytes(preHash))
	return
}

// AddBlock add a new block to blockchain
func (bc *BlockChain) AddBlock(block Block) bool {
	bc.Lock()
	defer bc.Unlock()
	// block varify successfully, add the block to the block chain
	if block.preHash.Equal(bc.blocks[bc.GetLength()-1].Hash()) && block.Varify() {
		bc.blocks = append(bc.blocks, block)
		bc.hashTable[block.Hash().String()] = &block
		bc.UnconfirmedData.Reset()
		return true
	}
	return false
}

// GetLength get the length of the block chain
func (bc *BlockChain) GetLength() int {
	return len(bc.blocks)
}

func (bc *BlockChain) GetBlockOfIndex(idx int) (*Block, bool) {
	if idx >= len(bc.blocks) {
		return nil, false
	}
	block := bc.blocks[idx]
	return &block, true
}

func (bc *BlockChain) GetBlockOfHash(hash string) (*Block, bool) {
	block, ok := bc.hashTable[hash]
	return block, ok
}

// GetBlocks get the blocks of the blockchain
func (bc *BlockChain) GetBlocks() []Block {
	return bc.blocks
}

func (bc *BlockChain) Lock() {
	bc.mutex.Lock()
}

func (bc *BlockChain) Unlock() {
	bc.mutex.Unlock()
}

func (bc *BlockChain) Update(bc2 *BlockChain) {
	bc.Lock()
	defer bc.Unlock()

	if bc2.GetLength() > bc.GetLength() {
		bc.blocks = bc2.blocks
		bc.hashTable = bc2.hashTable
	}
}
