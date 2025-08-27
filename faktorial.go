package main

import (
	"fmt"
	"math"
)

// fungsi faktorial
func factorial(n int) uint64 {
	if n == 0 {
		return 1
	}
	hasil := uint64(1)
	for i := 2; i <= n; i++ {
		hasil *= uint64(i)
	}
	return hasil
}

// fungsi f(n) = ceil(n! / 2^n)
func f(n int) uint64 {
	nomor := float64(factorial(n))
	den := math.Pow(2, float64(n))
	return uint64(math.Ceil(nomor / den))
}

func main() {
	for i := 0; i <= 10; i++ {
		fmt.Printf("f(%v) = %v \n", i, f(i))
	}
}
