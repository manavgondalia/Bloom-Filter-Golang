package hashes

import "github.com/spaolacci/murmur3"

// import "fmt"

const FNV_OFFSET_BASIS uint64 = 14695981039346656037
const FNV_PRIME uint64 = 1099511628211

func FNV_1(data string) uint64 {
	hash := FNV_OFFSET_BASIS
	for i := 0; i < len(data); i++ {
		hash *= FNV_PRIME
		hash ^= uint64(data[i])
	}
	return hash
}

func FNV_1A(data string) uint64 {
	hash := FNV_OFFSET_BASIS
	for i := 0; i < len(data); i++ {
		hash ^= uint64(data[i])
		hash *= FNV_PRIME
	}
	return hash
}

func FNV1A_VariantA(data string) uint64 {
	hash := FNV_OFFSET_BASIS
	for i := 0; i < len(data); i++ {
		hash ^= uint64(data[i])
		hash = (hash << 1) + (hash << 4) + (hash << 7) + (hash << 8) + (hash << 24)
	}
	return hash
}

func FNV1A_VariantB(data string) uint64 {

	hash := FNV_OFFSET_BASIS
	for i := 0; i < len(data); i++ {
		hash ^= uint64(data[i])
		hash = (hash << 1) ^ (hash >> 3) ^ (hash * FNV_PRIME)
	}
	return hash
}

func dbj2(data string) uint64 {
	hash := uint64(5381)
	for i := 0; i < len(data); i++ {
		hash = ((hash << 5) + hash) + uint64(data[i]) // 33
	}
	return hash
}

func HashFuncArrayGenerator(num_func uint64) []func(data string) uint64 {
	HashFunctionsArray := []func(data string) uint64{FNV_1, FNV_1A, FNV1A_VariantA, FNV1A_VariantB, dbj2}
	if num_func > 5 {
		// add more murmur3 with different seeds
		for i := 0; uint64(i) < num_func-5; i++ {
			func_to_add := func(data string) uint64 {
				seed := uint64(i)
				// implement murmur3 hash function
				return murmur3.Sum64WithSeed([]byte(data), uint32(seed))

			}
			HashFunctionsArray = append(HashFunctionsArray, func_to_add)
		}
	}
	return HashFunctionsArray
}

// var HashFunctionsArray = []func(data string) uint64{FNV_1, FNV_1A, FNV1A_VariantA, FNV1A_VariantB, dbj2}
