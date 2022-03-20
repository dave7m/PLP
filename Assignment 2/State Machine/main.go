package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	path string
)

// init is called before main()
func init() {
	// todo get default value for car.machine path or make it required
	flag.StringVar(&path, "path", "", "absolute path of machine file")
}

func main() {
	flag.Parse()

	if len(path) == 0 {
		fmt.Println("Usage: main.go -path")
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println("path: ", path)

	executeMachineAt(path)
}

func executeMachineAt(path string) {

}
