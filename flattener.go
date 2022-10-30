package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	baseDir                        string
	verbose, delete, overwrite     bool
	dirsProcessed, success, failed int
)

var (
	fBaseDir   = flag.String("dir", "", "<REQUIRED> Directory to be processed (use complete paths ex c:\\foo or /opt/foo)")
	fVerbose   = flag.Bool("verbose", false, "Verbose execution")
	fDelete    = flag.Bool("delete", false, "Delete processed sub folders")
	fOverwrite = flag.Bool("overwrite", false, "Overwrites files in destination if file with same name already exists")
)

func main() {
	flag.Parse()

	baseDir = strings.TrimSuffix(*fBaseDir, string(os.PathSeparator))
	verbose = *fVerbose
	delete = *fDelete
	overwrite = *fOverwrite

	if len(baseDir) == 0 {
		flag.Usage()
		os.Exit(-1)
	}
	isDir, err := isDirectory(baseDir)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	if !isDir {
		log.Fatal("Provided directory ", baseDir, " is not a directory")
		os.Exit(-1)
	}
	fmt.Println("Starting flattening process")
	processDir(baseDir)
	fmt.Println("Finished flattening process with:")
	fmt.Println("\t", dirsProcessed, " directories processed")
	fmt.Println("\t", success, " files successfully moved")
	fmt.Println("\t", failed, " files failed to moved")
}

func processDir(dirPath string) {
	dirsProcessed++
	if verbose {
		fmt.Println("Processing directory:", dirPath)
	}
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, f := range files {
		fullPath := dirPath + string(os.PathSeparator) + f.Name()
		destFullPath := baseDir + string(os.PathSeparator) + f.Name()
		if f.IsDir() {
			processDir(fullPath)
		} else {
			if baseDir == dirPath {
				continue
			}

			if verbose {
				fmt.Println("Moving file", fullPath, "to", baseDir)
			}
			if !overwrite {
				if _, err := os.Stat(destFullPath); !errors.Is(err, os.ErrNotExist) {
					fmt.Println("Cannot move file", fullPath, "to", destFullPath, "file with the same name already exists")
					failed++
					continue
				}
			}
			merr := os.Rename(fullPath, destFullPath)
			if merr != nil {
				failed++
				log.Fatal(merr)
			} else {
				success++
			}
		}
	}
	if delete && dirPath != baseDir {
		if verbose {
			fmt.Println("Removing directory", dirPath)
		}
		merr := os.Remove(dirPath)
		if merr != nil {
			log.Fatal(merr)
		}
	}

}

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
