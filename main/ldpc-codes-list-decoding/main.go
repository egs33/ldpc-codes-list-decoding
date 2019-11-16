package main

import (
	"../../pkg/channel"
	"../../pkg/ldpc"
	"fmt"
)

func main() {
	length := 1024
	bsc := channel.NewBinarySymmetricChannel(0.11)
	code := ldpc.ConstructCode(length, 512)
	codeword := make([]int, length)
	received := bsc.Channel(codeword)
	decoded := code.Decode(received)
	fmt.Println(decoded)
}
