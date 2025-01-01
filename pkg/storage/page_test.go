package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadPaeg(t *testing.T) {

	page := NewPage(30)
	page.PutInt(0, 1)
	page.PutInt(4, 2)
	page.PutInt(8, 3)
	page.PutInt(12, 4)
	page.PutInt(16, 5)
	page.PutInt(20, 6)

	assert.Equal(t, 1, page.GetInt(0))
	assert.Equal(t, 2, page.GetInt(4))
	assert.Equal(t, 3, page.GetInt(8))
	assert.Equal(t, 4, page.GetInt(12))
	assert.Equal(t, 5, page.GetInt(16))
	assert.Equal(t, 6, page.GetInt(20))

	page = NewPage(50)
	page.PutBytes(0, []byte("lintang")) /// 7 + 4
	page.PutBytes(11, []byte("birda"))  // 5 + 4
	page.PutBytes(20, []byte("saputra"))

	assert.Equal(t, "lintang", string(page.GetBytes(0)))
	assert.Equal(t, "birda", string(page.GetBytes(11)))
	assert.Equal(t, "saputra", string(page.GetBytes(20)))

	page = NewPage(50)
	page.putString(0, "lintang")
	page.putString(11, "birda")
	page.putString(20, "saputra")

	assert.Equal(t, "lintang", page.GetString(0))
	assert.Equal(t, "birda", page.GetString(11))
	assert.Equal(t, "saputra", page.GetString(20))

}
