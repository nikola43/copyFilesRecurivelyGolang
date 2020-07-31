package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/schollz/progressbar"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
	GreenColor   = "\033[0;32m%s\033[0m"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		fmt.Println("Uso: ./compressRecursively 'path/to/input/folder'")
		os.Exit(0)
	}

	root := os.Args[1]

	rootSlice := strings.Split(root, "/")
	copyPath := ""

	var i = 0
	for i = 0; i < len(rootSlice) - 1; i++ {
		copyPath += rootSlice[i]+"/"
	}
	copyPath += rootSlice[len(rootSlice) - 1]
	copyPath += "Copy"

	fmt.Println(copyPath)

	var successCounter = 0
	files, directories := getFilesAndDirectories(root)
	fmt.Printf(NoticeColor, "Número de ficheros: "+strconv.Itoa(len(files)))
	fmt.Printf(InfoColor, "\t Número de directorios: "+strconv.Itoa(len(directories)))
	fmt.Println("")

	FileExists(root)
	RemoveDirectory(copyPath)
	MakeDirectory(copyPath)

	// create and start new bar

	bar := progressbar.Default(int64(len(files)))
	fmt.Printf(WarningColor, "Comprimiendo...")
	fmt.Println("")
	var copyPath1 = rootSlice[len(rootSlice) -1] + "Copy"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path != root {
			var t = path
			t = strings.Replace(t, rootSlice[len(rootSlice) -1], copyPath1, -1)
			if info.IsDir() {
				directories = append(directories, path)
				MakeDirectory(t)
			} else {
				files = append(files, path)
				CopyFile(path, t)

				exist := FileExists(t)
				if exist {
					successCounter++
					bar.Add(1)
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf(GreenColor, "Número de ficheros copiados correctamente: "+strconv.Itoa(successCounter))
}

func getFilesAndDirectories(root string) ([]string, []string) {
	var files []string
	var directories []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			directories = append(directories, path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files, directories
}


// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func MakeDirectory(dirPath string) bool {
	err := os.Mkdir(dirPath, 0755)
	return checkError(err)
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println(err)
		//log.Fatal(err)
	}
	return true
}

func RemoveFile(filePath string) bool {
	err := os.Remove(filePath) // remove a single file
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func RemoveDirectory(dirPath string) bool {
	err := os.RemoveAll(dirPath) // Delete an entire directory
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return checkError(err)
}

func GenerateMD5HashFromFile(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	fmt.Println(returnMD5String)
	return returnMD5String, nil
}

/*
	hashIn, _ :=GenerateMD5HashFromFile(path)
	hashOut, _ :=GenerateMD5HashFromFile("assetsCopy/" + tempPath)
	if hashIn == hashOut {
		fmt.Println("OK")
	}
*/
/*
	fmt.Print("\033[G\033[K") // move the cursor left and clear the line
	if err != nil {
		fmt.Printf("Could not retrieve file info for %s\n", successCounter)
	} else {
		fmt.Printf("Retrieved %s\n", strconv.Itoa(successCounter))
	}
	fmt.Print("\033[A") // move the cursor up
*/
