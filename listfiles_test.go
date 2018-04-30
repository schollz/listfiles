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
	start := time.Now()
	dir1 := "."
	files, err := ListFiles(dir1)
	assert.Nil(t, err)
	smallNumFiles := len(files)

	dir2 := "../.."
	files, err = ListFiles(dir2)
	assert.Nil(t, err)
	largeNumFiles := len(files)

	fmt.Printf("\tNum files:\t%d\t\t%d\n", smallNumFiles, largeNumFiles)
	fmt.Print("ListFilesRecursively\t")
	start = time.Now()
	files, err = ListFilesRecursively(dir1)
	assert.Nil(t, err)
	fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
	start = time.Now()
	files, err = ListFilesRecursively(dir2)
	assert.Nil(t, err)
	fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())

	fmt.Print("ListFilesInParallel\t")
	start = time.Now()
	files, err = ListFilesRecursivelyInParallel(dir1)
	assert.Nil(t, err)
	fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
	start = time.Now()
	files, err = ListFilesRecursivelyInParallel(dir2)
	assert.Nil(t, err)
	fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())

	if runtime.GOOS != "windows" {
		fmt.Print("ListFilesUsingC\t")
		start = time.Now()
		files, err = ListFilesUsingC(dir1)
		assert.Nil(t, err)
		fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
		start = time.Now()
		files, err = ListFilesUsingC(dir2)
		assert.Nil(t, err)
		fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())
	}

	fmt.Print("ListFilesGodirwalk\t")
	start = time.Now()
	files, err = ListFilesGodirwalk(dir1)
	assert.Nil(t, err)
	fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
	start = time.Now()
	files, err = ListFilesGodirwalk(dir2)
	assert.Nil(t, err)
	fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())

}
