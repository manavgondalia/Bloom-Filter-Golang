package main

import (
	"encoding/binary"
	"fmt"
	"go-spell-checker/common"
	"log"
	"os"
	// Import the LoadBloomFilterFromFile function
)

type BloomFilter struct {
	BitSet        []bool
	K             uint64 // number of hash functions
	N             uint64 // number of elements
	M             uint64 // size of bitset
	HashFunctions []func(data string) uint64
}

func (bf *BloomFilter) Check(data string) bool {
	for _, hashFunc := range bf.HashFunctions {
		index := hashFunc(data) % uint64(len(bf.BitSet))
		if !bf.BitSet[index] {
			return false
		}
	}
	return true
}

// Load Bloom filter from file with custom header
func LoadBloomFilterFromFile(filename string) (*BloomFilter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read and parse custom header
	header := make([]byte, 12)
	_, err = file.Read(header)
	if err != nil {
		return nil, err
	}

	if string(header[0:4]) != "CCBF" {
		return nil, fmt.Errorf("invalid file format")
	}

	version := binary.BigEndian.Uint16(header[4:6])
	if version != 1 {
		return nil, fmt.Errorf("unsupported version: %d", version)
	}

	k := binary.BigEndian.Uint16(header[6:8])
	m := binary.BigEndian.Uint32(header[8:12])

	bf := &BloomFilter{
		BitSet:        make([]bool, m),
		K:             uint64(k),
		M:             uint64(m),
		N:             0, // This value is not stored in the file
		HashFunctions: common.HashFuncArrayGenerator(uint64(k))[:k],
	}

	// Read bitset
	bitSet := make([]byte, m)
	_, err = file.Read(bitSet)
	if err != nil {
		return nil, err
	}

	for i, b := range bitSet {
		if b == 1 {
			bf.BitSet[i] = true
		}
	}

	return bf, nil
}

func testing_accuracy(bf *BloomFilter, testWords []string) float64 {
	total := 0
	positive := 0

	for _, word := range testWords {
		total++
		if bf.Check(word) {
			positive++
		}
	}
	return float64(positive) / float64(total)
}

func main() {
	// Load the bloom filter from the file
	bf, err := LoadBloomFilterFromFile("words.bf")
	if err != nil {
		log.Fatalf("Error loading Bloom filter: %v", err)
	}

	// Test examples
	// testWords := []string{"hello", "world", "foo", "gondalia"}

	// for _, word := range testWords {
	// 	if bf.Check(word) {
	// 		fmt.Printf("Word '%s' is possibly in the Bloom filter.\n", word)
	// 	} else {
	// 		fmt.Printf("Word '%s' is definitely not in the Bloom filter.\n", word)
	// 	}
	// }

	// generate testwords from testing.txt
	// Read the file
	file, err := os.Open("testing.txt")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Read the words
	var word string
	var words []string
	for {
		_, err := fmt.Fscanf(file, "%s\n", &word)
		if err != nil {
			break
		}
		words = append(words, word)
	}

	// Test the accuracy
	accuracy := testing_accuracy(bf, words)
	fmt.Printf("Accuracy: %.2f%%\n", accuracy*100)

}
