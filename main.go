package main

import (
	"fmt"
	fileutils "github.com/nikola43/copyFilesRecurivelyGolang/utils"
	"github.com/schollz/progressbar"
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

	var ignore = false

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		fmt.Println("Usage: ./compressRecursively 'path/to/input/folder'")
		os.Exit(0)
	}

	root := os.Args[1]

	rootSlice := strings.Split(root, "/")
	copyPath := ""

	var i = 0
	for i = 0; i < len(rootSlice)-1; i++ {
		copyPath += rootSlice[i] + "/"
	}
	copyPath += rootSlice[len(rootSlice)-1]
	copyPath += "Copy"

	fmt.Println(copyPath)

	var successCounter = 0
	files, directories := fileutils.GetFilesAndDirectories(root)
	fmt.Printf(NoticeColor, "Files: "+strconv.Itoa(len(files)))
	fmt.Printf(InfoColor, "\t Directories: "+strconv.Itoa(len(directories)))
	fmt.Println("")

	fileutils.FileExists(root)
	fileutils.RemoveDirectory(copyPath)
	fileutils.MakeDirectory(copyPath)

	// create and start new bar

	bar := progressbar.Default(int64(len(files)))
	fmt.Printf(WarningColor, "Compressing...")
	fmt.Println("")
	var copyPath1 = rootSlice[len(rootSlice)-1] + "Copy"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path != root {
			var t = path
			t = strings.Replace(t, rootSlice[len(rootSlice)-1], copyPath1, -1)
			var compressedPath = t

			if info.IsDir() {
				//directories = append(directories, path)
				fileutils.MakeDirectory(t)
			} else {
				//files = append(files, path)

				// Get the content
				contentType, err := fileutils.GetFileContentType(path)
				if err != nil {
					panic(err)
				}

				//fmt.Println("Content Type: " + contentType)

				isAudio := strings.Contains(path, "mp3")

				if isAudio {
					fileutils.CopyFile(path, t)
				} else {
					if contentType == "image/jpeg" || contentType == "image/png" {

						isCompressed := strings.Contains(path, "compress")
						//fmt.Println(isCompressed) // true

						if isCompressed {
							fileutils.CopyFile(path, compressedPath)
						} else {
							err := fileutils.CompressImage(path, t, 50)
							if err != nil {
								fileutils.WriteLog("compress image " + path +" -> "+ err.Error())
							}
						}

					} else if contentType == "video/mp4" || contentType == "application/octet-stream" {
						isCompressed := strings.Contains(path, "compress")
						isJPG := strings.Contains(path, ".jpg")
						isPNG := strings.Contains(path, ".png")

						//fmt.Println(isCompressed) // true

						if isCompressed {
							fileutils.CopyFile(path, compressedPath)
						} else {
							fmt.Println(path)
							if isJPG || isPNG{
								ignore = true
								fileutils.CopyFile(path, t)
							} else {
								err := fileutils.CompressMP4(path, t)
								if err != nil {
									//panic(err)
									fileutils.WriteLog("compress video " + path +" -> "+ err.Error())
								}
							}
						}
					}
				}

				if ignore == false {
					exist := fileutils.FileExists(t)
					if exist {
						successCounter++
						err := bar.Add(1)
						if err != nil {
							panic(err)
						}
						fileutils.RemoveFile(path)
						fileutils.WriteLog("compress " + t + " -> " + "OK")
					}
				} else {
					successCounter++
					err := bar.Add(1)
					if err != nil {
						panic(err)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf(GreenColor, "Success files: "+strconv.Itoa(successCounter))
}
