package main

import (
	"log"
	"os"

	"Huffman/internal"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("got %d arguments, expected 3", len(os.Args)-1)
	}
	mode := os.Args[2]
	inputFileName := os.Args[3]
	outputFileName := os.Args[4]

	switch mode {
	case "encode":
		err := internal.Encode(inputFileName, outputFileName)
		if err != nil {
			log.Fatal(err)
		}
	case "decode":
		//TODO decode
	default:
		log.Fatalf(`invalid mode "%s", expected either "encode" or "decode"`, mode)
	}
}
