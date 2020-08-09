package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/h2non/bimg"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func ExecuteCommandVerbose(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	fmt.Println("Result: ")
	return err
}
func PrintSameLine(str string) {
	fmt.Print("\033[G\033[K") // move the cursor left and clear the line
	fmt.Printf("Retrieved %s\n", str)
	fmt.Print("\033[A") // move the cursor up
}

func CompressMP4(inFile string, outFile string) error {
	// check if input file exists
	err := CheckIfFileExists(inFile)
	if err != nil {
		return err
	}

	//RemoveFile(outFile)
	/*
		// check if input file exists
		err = CheckIfFileExists(outFile)
		if err != nil {
			return err
		} else {
			//if exists then remove
			fmt.Println("output file called " + outFile + " already exist. Removing...")
			removeError := os.Remove(outFile)
			if removeError != nil {
				return removeError
			}
			fmt.Println("file " + outFile + " has been removed successfully")
		}
	*/

	// extract audio from video using ffmpeg library
	// ffmpeg -i input.mp4 -vcodec h264 -acodec aac output.mp4
	err = ExecuteCommandVerbose("ffmpeg", "-i", inFile, "-vcodec", "h264", "-acodec", "mp3", outFile)
	if err != nil {
		return err
	}
	err = CheckIfFileExists(outFile)
	if err != nil {
		return err
	}
	return nil
}
func CompressImage(inputFilePath, outFilePath string, quality int) error {
	options := bimg.Options{
		Quality: quality,
	}

	// open file
	buffer, err := bimg.Read(inputFilePath)
	if err != nil {
		return err
	}

	// process file
	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		return err
	}

	// save file
	err = bimg.Write(outFilePath, newImage)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func WriteLog(message string) {
	date := fmt.Sprintf(time.Now().Format("2006-01-02"))
	var f *os.File
	var err error
	if fileExists("./logs/" + date + ".txt") {
		f, err = os.OpenFile("./logs/"+date+".txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		f, err = os.Create("./logs/" + date + ".txt")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	newLine := message + " | " + fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))

	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
}

/**
Check if file exist and its size
*/
func CheckIfFileExists(f string) error {
	var err error
	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(f); err == nil && fileInfo.Size() > 0 {
		return nil
	}
	return err
}
func GetFilesAndDirectories(root string) ([]string, []string) {
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
	return CheckError(err)
}

func MakeDirectory(dirPath string) bool {
	err := os.Mkdir(dirPath, 0755)
	return CheckError(err)
}

func CheckError(err error) bool {
	if err != nil {
		//fmt.Println(err)
		log.Fatal(err.Error())
	}
	return true
}

func RemoveFile(filePath string) bool {
	err := os.Remove(filePath) // remove a single file
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func GetFileContentType(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

/**
Add audio to video and generate video
*/
func AddAudioToVideo(inFile string, outFile string) (string, error) {
	// check if input file exists
	inFileError := CheckIfFileExists(inFile)
	if inFileError != nil {
		fmt.Printf("Error: %s", inFileError.Error())
		return "", inFileError
	}

	// extract audio from video using ffmpeg library
	cmd := exec.Command("ffmpeg", "-i", inFile, "-i", "music_video.m4a", "-c", "copy", "-shortest", outFile)
	//cmd := exec.Command("ffmpeg", "-i", inFile, "-i", "/home/ecodadys/go/src/github.com/nikola43/ecodadys_api/video_music.mp3", "-c", "copy", "-map", "0:v", "-map", "1:a", outFile)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error injectando audio " + err.Error())
		return "", err
	}

	return outFile, nil
}

/**
Extract audio from video and generate audio file
*/
func ExtractAudioFromVideo(inFile string, outFile string) (string, error) {
	// check if input file exists
	inFileError := CheckIfFileExists(inFile)
	if inFileError != nil {
		fmt.Printf("Error: %s", inFileError.Error())
		return "", inFileError
	}

	// check if output file exists
	outFileError := CheckIfFileExists(outFile)
	if outFileError != nil {
		fmt.Printf("Error: %s", outFileError.Error())
		return "", outFileError
	} else {
		// // if exists then remove
		// removeFileError := os.Remove(outFile)
		// if removeFileError != nil {
		// 	fmt.Printf("Error: %s", removeFileError.Error())
		// 	return removeFileError
		// }
	}

	// extract audio from video using ffmpeg library
	cmd := exec.Command("ffmpeg", "-i", inFile, "-f", "mp3", "-ab", "192000", "-vn", outFile)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return outFile, nil
}

func calculateCompressionPercentage(originalFileSize int64, outputFileSize int64) int64 {
	return (100 * outputFileSize) / originalFileSize
}

func printCompressionResult(inputFile, outputFile string) {
	fmt.Print(outputFile + "->" + strconv.FormatInt(getFileSize(outputFile), 10) + " " + strconv.FormatInt(100-(calculateCompressionPercentage(getFileSize(inputFile), getFileSize(outputFile))), 10) + "%")
}

func getFileSize(filePath string) int64 {
	file, err := os.Stat(filePath)
	CheckError(err)
	return file.Size()
}

/**
Extract audio from video and generate audio file
*/
func ExtractThumbnailFromVideo(inFile string, outFile string) error {
	// check if input file exists
	inFileError := CheckIfFileExists(inFile)
	if inFileError != nil {
		fmt.Printf("Error: %s", inFileError.Error())
		return inFileError
	}

	// check if output file exists
	outFileError := CheckIfFileExists(outFile)
	if outFileError != nil {
		fmt.Printf("Error: %s", outFileError.Error())
		return outFileError
	} else {
		// // if exists then remove
		// removeFileError := os.Remove(outFile)
		// if removeFileError != nil {
		// 	fmt.Printf("Error: %s", removeFileError.Error())
		// 	return removeFileError
		// }
	}

	// extract audio from video using ffmpeg library
	cmd := exec.Command("ffmpeg", "-i", inFile, "-ss", "00:00:05.000", "-vframes", "1", outFile)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
