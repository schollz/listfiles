package listfiles

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/mitchellh/hashstructure"
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
			Path:    path,
			Size:    f.Size(),
			Mode:    f.Mode(),
			ModTime: f.ModTime(),
			IsDir:   f.IsDir(),
		})
		return nil
	})
	return
}

// File is the object that contains the info and path of the file
type File struct {
	Path    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
	Hash    uint64 `hash:"ignore"`
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
			Path:    dir,
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		},
	}
	f.Close()
	fileChan := make(chan File)
	startedDirectories := make(chan bool)
	go listFilesInParallel(dir, startedDirectories, fileChan)

	runningCount := 1
	for {
		select {
		case file := <-fileChan:
			files = append(files, file)
		case newDir := <-startedDirectories:
			if newDir {
				runningCount++
			} else {
				runningCount--
			}
		default:
		}
		if runningCount == 0 {
			break
		}
	}
	return
}

func listFilesInParallel(dir string, startedDirectories chan bool, fileChan chan File) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fileStruct := File{
			Path:    path.Join(dir, f.Name()),
			Size:    f.Size(),
			Mode:    f.Mode(),
			ModTime: f.ModTime(),
			IsDir:   f.IsDir(),
		}
		h, err := hashstructure.Hash(fileStruct, nil)
		if err != nil {
			panic(err)
		}
		fileStruct.Hash = h
		fileChan <- fileStruct
		if f.IsDir() {
			startedDirectories <- true
			go listFilesInParallel(path.Join(dir, f.Name()), startedDirectories, fileChan)
		}
	}
	startedDirectories <- false
	return
}
