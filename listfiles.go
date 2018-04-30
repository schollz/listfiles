package listfiles

/*
#include <stdio.h>
#include <dirent.h>
#include <string.h>
#include <stdlib.h>
#include <limits.h>
#include <sys/stat.h>
extern void count(char *path, char *outfile);
*/
import "C"
import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/karrick/godirwalk"
	"github.com/mitchellh/hashstructure"
)

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func ListFiles(dir string) (files []File, err error) {
	if runtime.GOOS == "windows" {
		files, err = ListFilesRecursivelyInParallel(dir)
	} else {
		files, err = ListFilesUsingC(dir)
	}
	return
}

func ListFilesUsingC(dir string) (files []File, err error) {
	f, _ := ioutil.TempFile("/tmp", "listfiles")
	tempfile := f.Name()
	f.Close()
	defer os.Remove(tempfile)

	arg1 := C.CString(dir)
	defer C.free(unsafe.Pointer(arg1))
	arg2 := C.CString(tempfile)
	defer C.free(unsafe.Pointer(arg2))
	C.count(arg1, arg2)

	// count number of lines
	inFile, err := os.Open(tempfile)
	if err != nil {
		return
	}
	lines, err := lineCounter(inFile)
	inFile.Close()
	if err != nil {
		return
	}

	inFile, err = os.Open(tempfile)
	if err != nil {
		return
	}

	type result struct {
		err  error
		file File
	}
	jobs := make(chan string, lines)
	results := make(chan result, lines)

	for w := 0; w < runtime.NumCPU()*2; w++ {
		go func(jobs <-chan string, results chan<- result) {
			for path := range jobs {
				f, err := os.Lstat(path)
				if err != nil {
					results <- result{err: err}
				} else {
					file := File{
						Path:    path,
						Size:    f.Size(),
						Mode:    f.Mode(),
						ModTime: f.ModTime(),
						IsDir:   f.IsDir(),
					}
					h, err := hashstructure.Hash(file, nil)
					if err != nil {
						panic(err)
					}
					file.Hash = h

					results <- result{
						file: file,
						err:  nil,
					}
				}
			}
		}(jobs, results)
	}

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		path := filepath.Clean(strings.TrimSpace(scanner.Text()))
		jobs <- path
	}
	close(jobs)

	files = make([]File, lines)
	i := 0
	for j := 0; j < lines; j++ {
		result := <-results
		if result.err != nil {
			continue
		}
		files[i] = result.file
		i++
	}
	return
}

func ListFilesGodirwalk(dir string) (files []File, err error) {
	files = []File{}
	err = godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) (err error) {
			f, err := os.Stat(osPathname)
			if err != nil {
				return
			}
			file := File{
				Path:    osPathname,
				Size:    f.Size(),
				Mode:    f.Mode(),
				ModTime: f.ModTime(),
				IsDir:   f.IsDir(),
			}
			h, err := hashstructure.Hash(file, nil)
			if err != nil {
				return
			}
			file.Hash = h
			files = append(files, file)
			return nil
		},
		Unsorted: true,
		ScratchBuffer:  make([]byte, 64*1024),
	})
	return
}

// ListFilesRecursively uses filepath.Walk to list all the files
func ListFilesRecursively(dir string) (files []File, err error) {
	dir = filepath.Clean(dir)
	files = []File{}
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		file := File{
			Path:    path,
			Size:    f.Size(),
			Mode:    f.Mode(),
			ModTime: f.ModTime(),
			IsDir:   f.IsDir(),
		}
		h, err := hashstructure.Hash(file, nil)
		if err != nil {
			panic(err)
		}
		file.Hash = h
		files = append(files, file)
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
		{
			Path:    dir,
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		},
	}
	f.Close()

	h, err := hashstructure.Hash(files[0], nil)
	if err != nil {
		panic(err)
	}
	files[0].Hash = h

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
