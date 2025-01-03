package buffer

import (
	"fmt"
	"os"
	"testing"

	"github.com/lintang-b-s/go-simpledb/pkg/log"
	"github.com/lintang-b-s/go-simpledb/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func cleanDB() {
	stat, err := os.Stat("lintangdb")
	if err == nil && stat.IsDir() {
		os.RemoveAll("lintangdb")
	}
}

func TestBufferManager(t *testing.T) {
	cleanDB()

	dm := storage.NewDiskManager("lintangdb", 4096)
	lm, err := log.NewLogManager(dm, "lintangdb.log")
	if err != nil {
		t.Errorf("Error creating log manager: %s", err)
	}
	t.Run("success pin unpin buffer pool manager", func(t *testing.T) {

		bm := NewBufferPoolManager(5, dm, lm)
		buffers := make([]*Buffer, 10)
		blocks := make([]storage.BlockID, 10)
		for i := 0; i < 10; i++ {
			newBlock := storage.NewBlockID("test.db", i)
			blocks[i] = newBlock
		}

		for i := 0; i < 10; i++ {
			if i >= 5 {
				bm.UnpinPage(blocks[i-5], false)
			}
			bf, err := bm.PinPage(blocks[i])
			if err != nil {
				t.Errorf("Error pinning page: %s", err)
			}
			buffers[i] = bf
		}
	})

	t.Run("failed pin buffer pool manager because all buffer is pinned", func(t *testing.T) {

		bm := NewBufferPoolManager(5, dm, lm)

		buffers := make([]*Buffer, 10)
		blocks := make([]storage.BlockID, 10)
		for i := 0; i < 10; i++ {
			newBlock := storage.NewBlockID("test.db", i)
			blocks[i] = newBlock
		}

		for i := 0; i < 10; i++ {

			bf, err := bm.PinPage(blocks[i])
			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, fmt.Errorf("all pages are pinned"), err)
			}
			buffers[i] = bf
		}
	})

	t.Run("success fetch,create new, delete page - buffer pool manager", func(t *testing.T) {

		pages := make([]*storage.Page, 10)
		blocks := make([]storage.BlockID, 10)
		for i := 0; i < 10; i++ {

			page := storage.NewPage(4096)
			page.PutString(0, fmt.Sprintf("lintang%d", i))
			newBlockID := storage.NewBlockID("test.db", i)
			err = dm.Write(newBlockID, page)
			if err != nil {
				t.Error(err)
			}
			blocks[i] = newBlockID
		}

		bm := NewBufferPoolManager(5, dm, lm)

		// fetch all pages dan append ke pages
		for i := 0; i < 10; i++ {
			if i >= 5 {
				bm.UnpinPage(blocks[i-5], false)
			}
			bf, err := bm.FetchPage(blocks[i])
			if err != nil {
				t.Errorf("Error pinning page: %s", err)
			}
			pages[i] = bf
		}

		for idx, page := range pages {
			// read page
			block := blocks[idx]
			err = dm.Read(block, page)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, fmt.Sprintf("lintang%d", idx), page.GetString(0))
		}

	})
}
