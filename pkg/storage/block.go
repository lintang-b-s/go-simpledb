package storage

type BlockID struct {
	filename string
	blockNum int
}

func NewBlockID(filename string, blockNum int) *BlockID {

	return &BlockID{
		filename: filename,
		blockNum: blockNum,
	}
}

func (b *BlockID) getFilename() string {
	return b.filename
}

func (b *BlockID) getBlockNum() int {
	return b.blockNum
}
