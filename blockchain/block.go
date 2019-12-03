package blockchain

import (
	"bytes"
	// "fmt"
	"github.com/OhYee/cryptography_and_network_security/hash/sha"
	gb "github.com/OhYee/goutils/bytes"
	"math"
	"time"
)

type Block struct {
	timestamp int64
	blockData BlockData
	proof     uint64
	preHash   HashCode
}

// NewBlock init the block
func NewBlock(blockData BlockData, proof uint64, pre HashCode) Block {
	block := Block{
		timestamp: time.Now().Unix(),
		blockData: blockData.Copy(),
		proof:     proof,
		preHash:   pre,
	}
	return block
}

func (block *Block) GetTimestamp() int64 {
	return block.timestamp
}

func (block *Block) GetBlockData() BlockData {
	return block.blockData
}

func (block *Block) GetProof() uint64 {
	return block.proof
}

func (block *Block) GetPreHash() HashCode {
	return block.preHash
}

// Hash get the hash code of the block
func (block *Block) Hash() HashCode {
	return NewHashCodeFromBytes(sha.SHA256(block.ToBytes()))
}

// ToBytes block to []byte
func (block *Block) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(gb.FromInt64(block.timestamp))
	buf.Write([]byte{0x00, 0xFF})
	buf.Write(block.blockData.ToBytes())
	buf.Write([]byte{0x00, 0xFF})
	buf.Write(gb.FromUint64(block.proof))
	buf.Write([]byte{0x00, 0xFF})
	buf.Write(block.preHash.ToBytes())
	return buf.Bytes()
}

// Mine calc the hash
func (block Block) Mine() Block {
	for i := uint64(0); i < math.MaxUint64; i++ {
		block.proof = i
		if block.Varify() {
			break
		}
	}
	return block
}

// Varify a hash is varifity
func (block *Block) Varify() bool {
	return block.Hash().ToBytes()[0] == 0
}
