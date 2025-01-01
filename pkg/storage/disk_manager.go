package storage

import (
	"os"
)

type DiskManager struct {
	dbDir     string
	blockSize int
	isNew     bool
	openFiles map[string]*os.File
}

func NewDiskManager(dbDir string, blockSize int) *DiskManager {
	_, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		os.Mkdir(dbDir, 0755)
	}

	return &DiskManager{
		dbDir:     dbDir,
		blockSize: blockSize,
		isNew:     false,
		openFiles: make(map[string]*os.File),
	}
}

func (dm *DiskManager) Read(blockID *BlockID, page *Page) error {
	filename := dm.dbDir + "/" + blockID.getFilename()
	f, err := dm.getFile(filename)
	if err != nil {
		return err
	}
	_, err = f.Seek(int64(blockID.getBlockNum()*dm.blockSize), 0)
	if err != nil {
		return err
	}
	_, err = f.Read(page.contents())
	if err != nil {
		return err
	}
	return nil
}

func (dm *DiskManager) Write(blockID *BlockID, page *Page) error {
	filename := dm.dbDir + "/" + blockID.getFilename()
	f, err := dm.getFile(filename)
	if err != nil {
		return err
	}

	_, err = f.Seek(int64(blockID.getBlockNum()*dm.blockSize), 0)
	if err != nil {
		return err
	}

	_, err = f.Write(page.contents())
	if err != nil {
		return err
	}

	return nil
}

func (dm *DiskManager) Append(fileName string) (int, error) {
	newBlockNum, err := dm.blockLength(fileName)
	if err != nil {
		return 0, err
	}
	
	newBlock := NewBlockID(fileName, newBlockNum)
	b := make([]byte, dm.blockSize)
	f, err := dm.getFile(fileName)
	if err != nil {
		return 0, err
	}
	_, err = f.Seek(int64(newBlock.getBlockNum()*dm.blockSize), 0)
	if err != nil {
		return 0, err
	}
	_, err = f.Write(b)
	if err != nil {
		return 0, err
	}

	return newBlockNum, nil
}

func (dm *DiskManager) blockLength(fileName string) (int, error) {
	f, err := dm.getFile(fileName)
	if err != nil {
		return 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return int(fi.Size() / int64(dm.blockSize)), nil
}

func (dm *DiskManager) getFile(filename string) (*os.File, error) {
	file, exists := dm.openFiles[filename]
	var err error
	if !exists {
		file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0644)
		if err != nil {
			return nil, err
		}
	}

	dm.openFiles[filename] = file
	return file, nil
}
