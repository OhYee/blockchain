package blockchain

type BlockChain struct {
	UnconfirmedData BlockData         // the data wait for add to the blockchain
	blocks          []Block           // blocks of the blockchain
	hashTable       map[string]*Block // a map for hash to block
}

func NewBlockChain(initData BlockData) *BlockChain {
	initBlock := NewBlock(initData, 0, NewHashCodeFromBytes([]byte{}))
	initBlock = initBlock.Mine()

	initData.Reset()
	return &BlockChain{
		UnconfirmedData: initData,
		blocks:          []Block{initBlock},
		hashTable:       map[string]*Block{initBlock.Hash().String(): &initBlock},
	}
}

// ModifyData modify the data
func (bc *BlockChain) ModifyData(args ...interface{}) {
	bc.UnconfirmedData.Modify(args...)
}

// NewBlock init the block
func (bc *BlockChain) NewBlock() Block {
	return NewBlock(bc.UnconfirmedData, 0, bc.blocks[bc.GetLength()-1].Hash())
}

// AddBlock add a new block to blockchain
func (bc *BlockChain) AddBlock(block Block) {
	// block varify successfully, add the block to the block chain
	if block.Varify() {
		bc.blocks = append(bc.blocks, block)
		bc.hashTable[block.Hash().String()] = &block
		bc.UnconfirmedData.Reset()
	}
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
