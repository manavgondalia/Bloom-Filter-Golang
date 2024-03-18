package main

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



var HashFunctionsArray = []func(data string) uint64{FNV_1, FNV_1A, FNV1A_VariantA, FNV1A_VariantB, dbj2, }
