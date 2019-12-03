package blockchain

// BlockData the interface of the block data
type BlockData interface {
	Reset()                // reset the data to empty state
	Modify(...interface{}) // modify the data using the args
	ToBytes() []byte       // transfer to []byte
	Copy() BlockData       // make a copy of the data
}
