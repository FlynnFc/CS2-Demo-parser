package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var paths []string
var failed []string
var kills []events.Kill

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func main() {
	// ANSI escape codes for text colors

	allPlayerStats := []PlayerStats{}
	var root string

	//User enters path of folder
	fmt.Print("Please enter the path for the demo folder: ")
	fmt.Scan(&root)
	fmt.Println("Path=", root)

	err := filepath.WalkDir(root, visitFile)
	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", root, err)
	}

	fmt.Printf(Blue+"You are parsing %d files. Estimated time: %d - %d seconds\n"+Reset, len(paths), len(paths)-2, len(paths)+4)
	maxWorkers := 10
	overall := time.Now()
	pathChannel := make(chan string, len(paths))
	resultChannel := make(chan []PlayerStats, len(paths))

	for i := 0; i < maxWorkers; i++ {
		go func() {
			for path := range pathChannel {
				fmt.Printf(Yellow+"Parsing %s\n"+Reset, path)
				start := time.Now()
				game, err := demoParsing(path)
				checkError(err)
				elapsed := time.Since(start)
				fmt.Printf(Green+"Done! %s took %s\n"+Reset, path, elapsed)
				resultChannel <- game
			}
		}()
	}

	for _, path := range paths {
		pathChannel <- path
	}

	close(pathChannel)

	for i := 0; i < len(paths); i++ {
		allPlayerStats = append(allPlayerStats, <-resultChannel...)
	}

	fmt.Println(Green + "\nAll demos processed!" + Reset)
	fmt.Println()

	elapsed := time.Since(overall)

	if len(failed) > 1 {
		fmt.Printf(Red+"%d out of %d demos were invalid and were ignored.\n"+Reset, len(failed), len(paths))
	}
	fmt.Printf(Green+"Parsing %d demos took %s\n"+Reset, len(paths)-len(failed), elapsed)

	mergedData := combineStats(allPlayerStats)

	excelExporter(mergedData)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func visitFile(fp string, fi os.DirEntry, err error) error {
	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}
	if fi.IsDir() {
		return nil // not a file. ignore.
	}
	paths = append(paths, fp)
	return nil
}
