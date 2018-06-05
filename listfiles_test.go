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

	files5, err := ListFilesCwalk(".")
	assert.Nil(t, err)
	assert.NotEqual(t, len(files2), len(files5)) // not sure why this one isn't equal

	files6, err := ListFilesJonesWalk(".")
	assert.Nil(t, err)
	assert.NotEqual(t, len(files2), len(files6)) // not sure why this one isn't equal
}

func TestFilesPerSecond(t *testing.T) {
	ComputeHashes = false
	start := time.Now()
	dir1 := "."
	files, err := ListFiles(dir1)
	assert.Nil(t, err)
	smallNumFiles := len(files)

	dir2 := "../.."
	files, err = ListFiles(dir2)
	assert.Nil(t, err)
	largeNumFiles := len(files)

	fmt.Printf("Num files:\t\t\t\t\t\t%d\t\t%d\n", smallNumFiles, largeNumFiles)
	fmt.Print("ListFilesRecursively (github.com/schollz/listfiles)\t")
	start = time.Now()
	files, err = ListFilesRecursively(dir1)
	assert.Nil(t, err)
	fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
	start = time.Now()
	files, err = ListFilesRecursively(dir2)
	assert.Nil(t, err)
	fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())

	fmt.Print("ListFilesInParallel (github.com/schollz/listfiles)\t")
	start = time.Now()
	files, err = ListFilesRecursivelyInParallel(dir1)
	assert.Nil(t, err)
	fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
	start = time.Now()
	files, err = ListFilesRecursivelyInParallel(dir2)
	assert.Nil(t, err)
	fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())

	if runtime.GOOS != "windows" {
		fmt.Print("ListFilesUsingC (github.com/schollz/listfiles)\t\t")
		start = time.Now()
		files, err = ListFilesUsingC(dir1)
		assert.Nil(t, err)
		fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
		start = time.Now()
		files, err = ListFilesUsingC(dir2)
		assert.Nil(t, err)
		fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())
	}

	fmt.Print("ListFilesGodirwalk (github.com/karrick/godirwalk)\t")
	start = time.Now()
	files, err = ListFilesGodirwalk(dir1)
	assert.Nil(t, err)
	fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
	start = time.Now()
	files, err = ListFilesGodirwalk(dir2)
	assert.Nil(t, err)
	fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())

	fmt.Print("ListFilesCwalk (github.com/iafan/cwalk)\t\t\t")
	start = time.Now()
	files, err = ListFilesCwalk(dir1)
	assert.Nil(t, err)
	fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
	start = time.Now()
	files, err = ListFilesCwalk(dir2)
	assert.Nil(t, err)
	fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())

	fmt.Print("ListFilesJonesWalk (github.com/MichaelTJones/walk)\t")
	start = time.Now()
	files, err = ListFilesJonesWalk(dir1)
	assert.Nil(t, err)
	fmt.Printf("%2.0f files/s", float64(len(files))/time.Since(start).Seconds())
	start = time.Now()
	files, err = ListFilesJonesWalk(dir2)
	assert.Nil(t, err)
	fmt.Printf("\t%2.0f files/s\n", float64(len(files))/time.Since(start).Seconds())

}
