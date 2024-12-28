package internal

import (
	"Huffman/internal/models"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func restoreTree(inputFileName string) ([]models.Node, error) {
	log.Println("entered restoreTree()")
	defer log.Println("exiting restoreTree()")

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	log.Println("inputFile opened")

	scanner := bufio.NewScanner(inputFile)
	if !scanner.Scan() {
		return nil, io.EOF
	}
	treeSize, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, err
	}

	log.Printf("treeSize=%d\n", treeSize)

	tree := make([]models.Node, treeSize)
	for i := range treeSize {
		if !scanner.Scan() {
			return nil, io.EOF
		}
		var left, right, parent int
		var symbol byte
		_, err := fmt.Sscanf(scanner.Text(), "%d %d %d %d\n", &left, &right, &parent, &symbol)
		if err != nil {
			return nil, err
		}
		tree[i] = models.Node{Left: left, Right: right, Parent: parent, Symbol: symbol}
	}

	return tree, nil
}

func Unarchive(inputFileName, outputFileName string) error {
	tree, err := restoreTree(inputFileName)
	if err != nil {
		return err
	}
	for i, node := range tree {
		fmt.Printf("%d %d %d %d %d\n", i, node.Left, node.Right, node.Parent, node.Symbol)
	}
	// codes := generateCodes(tree)

	return nil
}
