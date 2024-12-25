package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var files []string

func main() {
	// read  fortune files
	fortuneCommand := exec.Command("fortune", "-f")
	// get the output of the command
	pipe, err := fortuneCommand.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}

	// start the command and read the output
	fortuneCommand.Start()
	outputStream := bufio.NewScanner(pipe)
	outputStream.Scan()

	// from the output, get the path of the fortune files
	line := outputStream.Text()
	path := line[strings.Index(line, "/"):]
	err = filepath.Walk(path, visitPathAndPopulateFiles)
	if err != nil {
		panic(err)
	}

	// get a random quote type and get a random quote from that type
	rand.Seed(time.Now().UnixNano())
	qouteTypeIndex := randomInt(1, len(files))
	randomFilePath := files[qouteTypeIndex]
	randomFile, err := os.Open(randomFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer randomFile.Close()

	// read the file and get a random quote
	b, err := ioutil.ReadAll(randomFile)
	if err != nil {
		panic(err)
	}
	quotes := string(b)
	quotesSlice := strings.Split(quotes, "%")
	qouteindex := randomInt(1, len(quotesSlice))
	fmt.Println(quotesSlice[qouteindex])

}

func visitPathAndPopulateFiles(path string, f os.FileInfo, err error) error {
	if strings.Contains(path, "/off/") {
		return nil
	}
	if filepath.Ext(path) == ".dat" {
		return nil
	}
	if f.IsDir() {
		return nil
	}
	files = append(files, path)
	return nil
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
