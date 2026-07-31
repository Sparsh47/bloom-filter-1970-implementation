package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Sparsh47/bloom-filter-implementation/utils"
)

const SIZE = 500000

const HASH_COUNT = 5

func insertWord(bitArray []bool, word string) {

	for i := range HASH_COUNT {
		h := utils.Hash([]byte(word), i*10)
		bitArray[h%SIZE] = true
	}
}

func queryWord(bitArray []bool, word string) (string, bool) {
	for i := range HASH_COUNT {
		h := utils.Hash([]byte(word), i*10)
		if !bitArray[h%SIZE] {
			return "Definitely Absent", false
		}
	}

	return "Possibly Present", true
}

func loadWords(file *os.File) ([]string, []string) {
	var wordsArray []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if word == "" {
			continue
		}

		if !utils.IsAlpha(word) {
			continue
		}

		wordsArray = append(wordsArray, word)

		if len(wordsArray) == 120000 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return wordsArray[:100000], wordsArray[100000:]
}

func main() {
	var false_positives int
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	bitArray := make([]bool, SIZE)

	words, dictWords := loadWords(file)

	for _, word := range words {
		insertWord(bitArray, word)
	}

	for _, word := range dictWords {
		_, result := queryWord(bitArray, word)

		if result {
			false_positives++
		}
	}

	measured_rate := float64(false_positives) / float64(len(dictWords))

	fmt.Printf("Loaded %d words\n", len(words))

	fmt.Printf("Loaded %d dict words\n", len(dictWords))

	fmt.Printf("Bit array length: %d\n", len(bitArray))

	fmt.Printf("False positives: %d\n", false_positives)

	fmt.Printf("Measured rate: %.4f\n", measured_rate)
}
