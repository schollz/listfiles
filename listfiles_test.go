package listfiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ListFilesRecursively("../croc")
	}
}

func BenchmarkListInParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ListFilesRecursivelyInParallel("../croc")
	}
}

func TestListFiles(t *testing.T) {
	files, err := ListFilesRecursively(".")
	assert.Nil(t, err)
	assert.Equal(t, 29, len(files))

	files, err = ListFilesRecursivelyInParallel(".")
	assert.Nil(t, err)
	assert.Equal(t, 29, len(files))
}
