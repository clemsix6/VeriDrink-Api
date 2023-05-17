package main

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"sync"
)

// String array to store the lines of the file.
var lines []string
var mu sync.Mutex

// Load the file and store lines in the global variable.
func loadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.New("failed to open the file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return errors.New("error reading the file")
	}

	return nil
}

// Get a random line from the global variable.
func getRandomLine() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if len(lines) == 0 {
		return "", errors.New("no lines in the file")
	}

	return lines[rand.Intn(len(lines))], nil
}
