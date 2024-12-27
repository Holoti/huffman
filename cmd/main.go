package main

import (
	"Huffman/internal"
	"Huffman/internal/constants"
	"log"
)

func main() {
	err := internal.Encode(constants.InputFileName, constants.OutputFileName)
	if err != nil {
		log.Fatal(err)
	}
}
