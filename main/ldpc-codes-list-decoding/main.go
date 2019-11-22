package main

import (
	"fmt"
	"github.com/egs33/ldpc-codes-list-decoding/pkg/channel"
	"github.com/egs33/ldpc-codes-list-decoding/pkg/ldpc"
)

func main() {
	length := 1024
	bsc := channel.NewBinarySymmetricChannel(0.11)
	code, _ := ldpc.ConstructCode(length, 512, 3, 6)
	codeword := make([]int, length)
	received := bsc.Channel(codeword)
	decoded := code.Decode(received)
	fmt.Println(decoded)
}
