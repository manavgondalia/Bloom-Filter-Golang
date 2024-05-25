package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"go-spell-checker/hashes"
	"math"
	"os"
)

type BloomFilter struct {
	BitSet        []bool
	K             uint64 // number of hash functions
	N             uint64 // number of elements
	M             uint64 // size of bitset
	HashFunctions []func(data string) uint64
}

func NewBloomFilter(k uint64, m uint64, n uint64) *BloomFilter {
	return &BloomFilter{
		BitSet:        make([]bool, m),
		K:             k,
		M:             m,
		N:             n,
		HashFunctions: hashes.HashFuncArrayGenerator(k)[:k],
		// if more hash fucntions needed, use murmur3 hash functions with different seeds
	}
}

func (bf *BloomFilter) Add(data string) {
	for _, hashFunc := range bf.HashFunctions {
		index := hashFunc(data) % uint64(len(bf.BitSet))
		bf.BitSet[index] = true
	}
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

// Save Bloom filter to file with custom header
func (bf *BloomFilter) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write custom header
	header := make([]byte, 12)
	copy(header[0:4], "CCBF")
	binary.BigEndian.PutUint16(header[4:6], 1)             // version number
	binary.BigEndian.PutUint16(header[6:8], uint16(bf.K))  // number of hash functions
	binary.BigEndian.PutUint32(header[8:12], uint32(bf.M)) // number of bits

	_, err = file.Write(header)
	if err != nil {
		return err
	}

	// Write bitset
	for _, bit := range bf.BitSet {
		var b byte
		if bit {
			b = 1
		}
		_, err := file.Write([]byte{b})
		if err != nil {
			return err
		}
	}

	return nil
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
		HashFunctions: hashes.HashFuncArrayGenerator(uint64(k))[:k],
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

func LoadAndUse() {
	bf, err := LoadBloomFilterFromFile("words.bf")
	if err != nil {
		fmt.Println("Error loading Bloom filter:", err)
		return
	}

	// Can modify this as per the scenario of checking membership
	testWords := []string{"hello", "world", "foo", "BAR"}

	for _, word := range testWords {
		if bf.Check(word) {
			fmt.Printf("Word '%s' is possibly in the Bloom filter.\n", word)
		} else {
			fmt.Printf("Word '%s' is definitely not in the Bloom filter.\n", word)
		}
	}
}

func main() {

	// Take false probability rate and number of elements from command line arguments
	load := flag.Int("load", 0, "Whether to load Bloom filter from file")
	build := flag.String("build", "words.txt", "Path to dictionary file to build Bloom filter")
	fpRate := flag.Float64("fp", 0.01, "False probability rate")
	numElements := flag.Int("n", -1, "Number of elements")

	flag.Parse()

	if *load == 1 {
		LoadAndUse()
		return
	}

	if *build == "" {
		fmt.Println("Please specify the path to the dictionary file using the -build flag.")
		return
	}

	file, err := os.Open(*build)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	var word string

	if *numElements == -1 {
		*numElements = 0
		for {
			_, err := fmt.Fscanln(file, &word)
			if err != nil {
				break
			}
			*numElements++
		}
	}

	// Calculate size of bitset and number of hash functions
	size := -1 * float64(*numElements) * (math.Log(*fpRate) / math.Pow(math.Log(2), 2))
	size = math.Ceil(size)

	k := (size / float64(*numElements)) * math.Log(2)
	k = math.Ceil(k)

	bf := NewBloomFilter(uint64(k), uint64(size), uint64(*numElements))

	fmt.Println("Size of bitset: ", bf.M)
	fmt.Println("Number of hash functions: ", bf.K)
	fmt.Println("Number of elements: ", bf.N)

	// Add elements from the dictionary file to the Bloom filter
	file, err = os.Open(*build)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	for i := 0; i < *numElements; i++ {
		_, err := fmt.Fscanln(file, &word)
		if err != nil {
			break
		}
		bf.Add(word)
	}

	// Dump the bloom filter to the file to disk for loading later

	err = bf.SaveToFile("words.bf")
	if err != nil {
		fmt.Println("Error saving Bloom filter:", err)
		return
	}

	fmt.Println("Bloom filter saved to words.bf")

	// Test examples
	fmt.Println(bf.Check("hello"))
	fmt.Println(bf.Check("world"))
	fmt.Println(bf.Check("foo"))
}
