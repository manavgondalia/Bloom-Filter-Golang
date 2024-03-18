package main

import (
	"flag"
	"fmt"
	// "fmt"
	"math"
)

type BloomFilter struct {
	bitSet        []bool
	k             uint64 // number of hash functions
	n             uint64 // number of elements
	m             uint64 // size of bitset
	hashFunctions []func(data string) uint64
}

func NewBloomFilter(k uint64, m uint64, n uint64) *BloomFilter {
	return &BloomFilter{
		bitSet:        make([]bool, m),
		k:             k,
		m:             m,
		n:             n,
		hashFunctions: HashFunctionsArray[:k],
	}
}

func (bf *BloomFilter) Add(data string) {
	for _, hashFunc := range bf.hashFunctions {
		index := hashFunc(data) % uint64(len(bf.bitSet))
		bf.bitSet[index] = true
	}
}

func (bf *BloomFilter) Check(data string) bool {
	for _, hashFunc := range bf.hashFunctions {
		index := hashFunc(data) % uint64(len(bf.bitSet))
		if !bf.bitSet[index] {
			return false
		}
	}
	return true
}

func main() {

	// take false probability rate and number of elements from command line arguments

	fpRate := flag.Float64("fp", 0.01, "False probability rate")
	numElements := flag.Int("n", 1000, "Number of elements")

	flag.Parse()

	// calculate size of bitset and number of hash functions

	size := -1 * float64(*numElements) * (math.Log(*fpRate) / math.Pow(math.Log(2), 2))
	size = math.Ceil(size)

	k := (size / float64(*numElements)) * math.Log(2)
	k = math.Ceil(k)

	bf := NewBloomFilter(uint64(k), uint64(size), uint64(*numElements))

	fmt.Println("Size of bitset: ", bf.m)
	fmt.Println("Number of hash functions: ", bf.k)

	bf.Add("hello")
	bf.Add("world")

	fmt.Println(bf.Check("hello"))
	fmt.Println(bf.Check("world"))
	fmt.Println(bf.Check("foo"))

}
