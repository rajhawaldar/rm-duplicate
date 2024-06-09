package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

var (
	dirPath     = flag.String("p", "", "Path to look for duplicate files")
	dryRun      = flag.Bool("d", false, "Dry Run(Print the duplicate files, DO NOT DELETE)")
	deleteCount = 0
)

type File struct {
	path   string
	remove bool
}

var count int

type Record map[string][]File

var exclude []string

func checkExcludedPathExist(path string) bool {
	for _, word := range exclude {
		if strings.Contains(path, word) {
			return true
		}
	}
	return false
}
func main() {
	flag.Parse()
	fileRecords := make(Record)
	exclude = append(exclude, ".git", ".vscode", ".mp4")
	if len(*dirPath) == 0 {
		fmt.Println("Please use --help for how to use rm-duplicate.")
		return
	}

	err := filepath.Walk(*dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Panic(err.Error())
			return err
		}
		if !info.IsDir() && !checkExcludedPathExist(path) {
			var fs = afero.NewOsFs()
			file, err := fs.Open(path)
			if err != nil {
				log.Panic("failed to open a file", path)
			}
			hash := md5.New()

			if _, err := io.Copy(hash, file); err != nil {
				log.Panic("calculating file hash failed")
			}
			file.Close()
			hashSum := string(hash.Sum(nil))
			fileRecords[hashSum] = append(fileRecords[hashSum], File{
				path:   path,
				remove: false,
			})
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	for key, records := range fileRecords {
		if len(records) == 1 {
			delete(fileRecords, key)
		}
	}
	for _, records := range fileRecords {
		count += len(records)
		fmt.Println("Length: ", len(records))
		for index, record := range records {
			if index == 0 {
				fmt.Println("Skipped: ", record.path)
				continue
			}
			if !*dryRun {
				// os.Remove(record.path)
			}
			fmt.Println("Removing: ", record.path)
		}
	}
	fmt.Println("Total", count, "files deleted!")
}
