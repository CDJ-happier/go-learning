package baseencode

import (
	"errors"
	"fmt"
)

const (
	carry    = 58
	maxBytes = 8
)

var (
	// ErrInvalidInput base error for invalid input
	ErrInvalidInput = errors.New("base58 invalid input")
	// ErrBase58Overflow is returned when the number to decode is too large, greater than eight bytes
	ErrBase58Overflow = fmt.Errorf("%w: number is too large", ErrInvalidInput)
	// ErrorInvalidCharacter is returned when an invalid character is found in the input
	ErrorInvalidCharacter = fmt.Errorf("%w: invalid character in input", ErrInvalidInput)

	chars = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

	// indices is an array of indices of characters in the base58 alphabet
	indices = [256]int{}

	// pow58 returns 58^n
	pow58 = [9]int{}
)

func init() {
	for i := range indices {
		indices[i] = -1
	}

	for i, char := range chars {
		indices[char] = i
	}

	for i := range pow58 {
		pow58[i] = pow(carry, i)
	}
}

// Base58Encode encodes a number to base58
func Base58Encode(num uint64) []byte {
	b := make([]byte, 0, maxBytes) // 58^8 = 1.28e14, enough for our use case

	for ; num > 0; num /= carry {
		b = append(b, chars[num%carry])
	}
	reverse(b)

	return b
}

// Base58Decode decodes a base58 encoded string to a number
func Base58Decode(b []byte) (uint64, error) {
	n := len(b)
	if n > maxBytes { // 58^8 = 1.28e14, less than math.MaxUint64, and it's enough for our use case
		return 0, ErrBase58Overflow
	}

	var num uint64

	for i := range n {
		pos := indices[b[i]]
		if pos == -1 {
			return 0, ErrorInvalidCharacter
		}

		num += uint64(pow58[n-i-1] * pos)
	}

	return num, nil
}

func reverse(a []byte) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

func pow(x, n int) int {
	if n == 0 {
		return 1
	}

	if n%2 == 0 {
		return pow(x*x, n/2) //nolint:mnd
	}

	return x * pow(x, n-1)
}
