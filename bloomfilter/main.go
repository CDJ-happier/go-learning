package main

import (
	"fmt"
	"hash/fnv"
	"math"
)

// Bloom Filter
// n: number of items 需要插入的元素数量
// p: error rate 允许的错误率
// m: bit array size
// k: number of hash functions
// bloom filter一般提供n, p. 并计算出相应的m, k.
// m = -n*ln(p) / (ln2)^2
// k = m/n * ln2

type BloomFilter struct {
	m      int
	k      int
	bitset []bool
}

// NewBloomFilter creates a new Bloom filter with m bits and k hash functions.
func NewBloomFilter(m int, k int) *BloomFilter {
	return &BloomFilter{m: m, k: k, bitset: make([]bool, m)}
}

// NewBloomFilter2 creates a new Bloom filter with n items and error rate p.
func NewBloomFilter2(n int, p float64) *BloomFilter {
	m := -float64(n) * math.Log(p) / (math.Ln2 * math.Ln2)
	k := math.Ln2 * m / float64(n)
	return NewBloomFilter(int(m)+1, int(math.Ceil(k)))
}

func (bf *BloomFilter) Add(item string) {
	for i := 0; i < bf.k; i++ {
		idx := bf.hashing([]byte(item), i)
		bf.bitset[idx] = true
	}
}

func (bf *BloomFilter) Contains(item string) bool {
	for i := 0; i < bf.k; i++ {
		idx := bf.hashing([]byte(item), i)
		if !bf.bitset[idx] {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) Reset() {
	for i := 0; i < bf.m; i++ {
		bf.bitset[i] = false
	}
}

func (bf *BloomFilter) hashing(data []byte, i int) uint64 {
	hasher := fnv.New64a()
	hasher.Write(data)
	hv := hasher.Sum64()
	hv += uint64(i) * 101 // 添加偏移
	return hv % uint64(bf.m)
}

func main() {
	n := 300
	truePositive := 0
	bloomFilter := NewBloomFilter2(n, 0.1)
	fmt.Printf("bloom filter size: %d, %d\n", bloomFilter.m, bloomFilter.k)
	for i := 0; i < n; i += 3 {
		truePositive++
		bloomFilter.Add(fmt.Sprintf("key%d", i))
	}
	positive := 0
	for i := 0; i < n; i++ {
		if bloomFilter.Contains(fmt.Sprintf("key%d", i)) {
			positive++
			fmt.Printf("key%d: maybe exist\n", i)
		}
	}
	fmt.Printf("true positive: %d, positive: %d\n", truePositive, positive)
	fmt.Printf("accuracy = %v\n", 1-math.Abs(float64(truePositive)-float64(positive))/float64(n))
}
