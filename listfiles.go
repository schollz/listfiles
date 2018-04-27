package listfiles

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func main() {
	fmt.Println(ListFilesRecursively("."))
	fmt.Println(ListFilesRecursivelyInParallel("."))
}

// ListFilesRecursively uses filepath.Walk to list all the files
func ListFilesRecursively(dir string) (files []File, err error) {
	dir = filepath.Clean(dir)
	files = []File{}
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		files = append(files, File{
			Info: f,
			Path: path,
		})
		return nil
	})
	return
}

// File is the object that contains the info and path of the file
type File struct {
	Info os.FileInfo
	Path string
}

// ListFilesRecursivelyInParallel uses goroutines to list all the files
func ListFilesRecursivelyInParallel(dir string) (files []File, err error) {
	dir = filepath.Clean(dir)
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	info, err := f.Stat()
	if err != nil {
		return
	}
	files = []File{
		File{
			Path: dir,
			Info: info,
		},
	}
	f.Close()
	fileChan := make(chan File)
	startedDirectories := make(chan bool)
	finishedDirectories := make(chan bool)
	go listFilesInParallel(dir, startedDirectories, finishedDirectories, fileChan)

	startedNum := 1
	finishedNum := 0
	for {
		select {
		case file := <-fileChan:
			files = append(files, file)
		case <-startedDirectories:
			startedNum += 1
		case <-finishedDirectories:
			finishedNum += 1
		default:
		}
		if startedNum == finishedNum {
			break
		}
	}
	return
}

func listFilesInParallel(dir string, startedDirectories chan bool, finishedDirectories chan bool, fileChan chan File) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fileChan <- File{
			Path: path.Join(dir, f.Name()),
			Info: f,
		}
		if f.IsDir() {
			startedDirectories <- true
			go listFilesInParallel(path.Join(dir, f.Name()), startedDirectories, finishedDirectories, fileChan)
		}
	}
	finishedDirectories <- true
	return
}
