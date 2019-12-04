package blockchain

// BlockData the interface of the block data
type BlockData interface {
	Reset()                                // reset the data to empty state
	Modify(args ...interface{})            // modify the data using the args
	ToBytes() []byte                       // transfer to []byte
	Verify() bool                          // verify the data is valid
	Copy() BlockData                       // make a copy of the data
	String(prefix string) string           // return the string of the data
	FromBytes(b []byte) (BlockData, error) // from []byte to generate the block data
}
