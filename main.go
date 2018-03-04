package main

import (
	"./Gonome"
	"fmt"
	"os"
	"strconv"
)

// Database location goes here
var DATABASE = "./Data/ProteinDatabase.txt"
var INDEX int = 0

func main() {

	// You can get individual args with normal indexing.
	file_name := os.Args[1]

	var arg_len = len(os.Args)
	if arg_len > 2 {
		DATABASE = os.Args[2]
		if arg_len > 3 {
			i, _ := strconv.Atoi(os.Args[3])
			INDEX = i
		}
	}

	fmt.Println("Input file: ", file_name)
	Gonome.StartSearch(file_name, DATABASE, INDEX)

}
