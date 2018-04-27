package listfiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ListFilesRecursively("../../")
	}
}

func BenchmarkListInParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ListFilesRecursivelyInParallel("../../")
	}
}

func BenchmarkListFromFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ListFilesFromFile("../../")
	}
}

func TestListFiles(t *testing.T) {
	files1, err := ListFilesRecursively(".")
	assert.Nil(t, err)

	files2, err := ListFilesRecursivelyInParallel(".")
	assert.Nil(t, err)
	assert.Equal(t, len(files2), len(files1))

	// files3, err := ListFilesFromFile(".")
	// assert.Nil(t, err)
	// assert.Equal(t, len(files2), len(files3))
}
