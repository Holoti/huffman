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
	mode := os.Args[1]
	inputFileName := os.Args[2]
	outputFileName := os.Args[3]

	switch mode {
	case "encode":
		if err := internal.Archive(inputFileName, outputFileName); err != nil {
			log.Fatal(err)
		}
	case "decode":
		//TODO decode
		if err := internal.Unarchive(inputFileName, outputFileName); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf(`invalid mode "%s", expected either "encode" or "decode"`, mode)
	}
}

// go run ./cmd/main.go "encode" "data/original.txt" "data/encoded.txt"
// go run ./cmd/main.go "decode" "data/encoded.txt" "data/decoded.txt"
