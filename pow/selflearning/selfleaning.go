package main

import (
	"math/big"
	"fmt"
)

const targetBits = 24

func main() {
	target := big.NewInt(1)
	fmt.Printf("target = %v\n", target)
	target.Lsh(target, uint(256-targetBits))
	fmt.Printf("target = %v\n", target)
}
