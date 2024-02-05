package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	maxStations := defaultStations
	maxRows := defaultRows

	nArgs := len(os.Args)

	if nArgs >= 3 {
		maxStations, _ = strconv.Atoi(os.Args[2])
	}

	if nArgs >= 2 {
		maxRows, _ = strconv.Atoi(os.Args[1])
	}

	file, _ := os.Open(stationsFile)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	fmt.Println(maxStations)
	fmt.Println(maxRows)
}
