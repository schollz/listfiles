package listfiles

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkListStdLib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ListFilesRecursively("../../")
	}
}

func BenchmarkListInParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ListFilesRecursivelyInParallel("../../")
	}
}

func BenchmarkListFilesUsingC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ListFilesUsingC("../../")
	}
}

func TestListFiles(t *testing.T) {
	files1, err := ListFilesRecursively(".")
	assert.Nil(t, err)

	files2, err := ListFilesRecursivelyInParallel(".")
	assert.Nil(t, err)
	assert.Equal(t, len(files2), len(files1))

	files3, err := ListFilesUsingC(".")
	assert.Nil(t, err)
	assert.Equal(t, len(files2), len(files3)+1)

	files4, err := ListFilesGodirwalk(".")
	assert.Nil(t, err)
	assert.Equal(t, len(files2), len(files4))
}

func TestFilesPerSecond(t *testing.T) {
	dir := "."
	start := time.Now()
	files, err := ListFilesRecursively(dir)
	assert.Nil(t, err)
	fmt.Printf("ListFilesRecursively %2.0f files/s (%d files)\n", float64(len(files))/time.Since(start).Seconds(), len(files))

	start = time.Now()
	files, err = ListFilesRecursivelyInParallel(dir)
	assert.Nil(t, err)
	fmt.Printf("ListFilesRecursivelyInParallel %2.0f files/s (%d files)\n", float64(len(files))/time.Since(start).Seconds(), len(files))

	if runtime.GOOS != "windows" {
		start = time.Now()
		files, err = ListFilesUsingC(dir)
		assert.Nil(t, err)
		fmt.Printf("ListFilesUsingC %2.0f files/s (%d files)\n", float64(len(files))/time.Since(start).Seconds(), len(files))
	}

	start = time.Now()
	files, err = ListFilesGodirwalk(dir)
	assert.Nil(t, err)
	fmt.Printf("ListFilesGodirwalk %2.0f files/s (%d files)\n", float64(len(files))/time.Since(start).Seconds(), len(files))

	dir = "../../"
	start = time.Now()
	files, err = ListFilesRecursively(dir)
	assert.Nil(t, err)
	fmt.Printf("ListFilesRecursively %2.0f files/s (%d files)\n", float64(len(files))/time.Since(start).Seconds(), len(files))

	start = time.Now()
	files, err = ListFilesRecursivelyInParallel(dir)
	assert.Nil(t, err)
	fmt.Printf("ListFilesRecursivelyInParallel %2.0f files/s (%d files)\n", float64(len(files))/time.Since(start).Seconds(), len(files))

	if runtime.GOOS != "windows" {
		start = time.Now()
		files, err = ListFilesUsingC(dir)
		assert.Nil(t, err)
		fmt.Printf("ListFilesUsingC %2.0f files/s (%d files)\n", float64(len(files))/time.Since(start).Seconds(), len(files))
	}

	start = time.Now()
	files, err = ListFilesGodirwalk(dir)
	assert.Nil(t, err)
	fmt.Printf("ListFilesGodirwalk %2.0f files/s (%d files)\n", float64(len(files))/time.Since(start).Seconds(), len(files))

}
