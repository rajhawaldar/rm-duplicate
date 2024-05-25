package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var (
	dirPath     = flag.String("p", "", "Path to look for duplicate files")
	deleteCount = 0
)

type DuplicateEntry map[string][]string

func main() {
	flag.Parse()
	dList := make(DuplicateEntry)
	if len(*dirPath) == 0 {
		fmt.Println("Please use --help for how to use rm-duplicate.")
		return
	}

	err := filepath.Walk(*dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Panic("failed to read directory.")
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				log.Panic("failed to open a file", path)
			}
			defer file.Close()
			hash := sha256.New()

			if _, err := io.Copy(hash, file); err != nil {
				log.Panic("calculating file hash failed")
			}
			hashSum := string(hash.Sum(nil))
			dList[hashSum] = append(dList[hashSum], file.Name())
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	for _, list := range dList {
		if len(list) > 1 {
			for _, file := range list[1:] {
				os.Remove(file)
				fmt.Println(file, "deleted!")
				deleteCount++
			}
		}
	}
	fmt.Println("Total", deleteCount, "files deleted!")
}
