package main

type BloomFilter struct {
	bitSet        []bool
	hashFunctions []func(data string) uint
}

func NewBloomFilter(size int, hashFuncs ...func(data string) uint) *BloomFilter {
	return &BloomFilter{
		bitSet:        make([]bool, size),
		hashFunctions: hashFuncs,
	}
}

func (bf *BloomFilter) Add(data string) {
	for _, hashFunc := range bf.hashFunctions {
		index := hashFunc(data) % uint(len(bf.bitSet))
		bf.bitSet[index] = true
	}
}

func (bf *BloomFilter) Check(data string) bool {
	for _, hashFunc := range bf.hashFunctions {
		index := hashFunc(data) % uint(len(bf.bitSet))
		if !bf.bitSet[index] {
			return false
		}
	}
	return true
}

func main() {
	bf := NewBloomFilter(1000, hashFunc1, hashFunc2)
	bf.Add("hello")
	bf.Add("world")

	println(bf.Check("hello")) 
	println(bf.Check("world")) 
	println(bf.Check("foo"))   
}

func hashFunc1(data string) uint {
	hash := uint(0)
	for i := 0; i < len(data); i++ {
		hash += uint(data[i])
	}
	return hash
}

func hashFunc2(data string) uint {
	hash := uint(0)
	for i := 0; i < len(data); i++ {
		hash += uint(data[i]) * 31
	}
	return hash
}
