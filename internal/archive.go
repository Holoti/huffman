package internal

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"Huffman/internal/models"
)

func createForest(inputFileName string) ([]models.Forest, []byte, error) {
	file, err := os.Open(inputFileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	frequency := make(map[byte]int)
	current_symbol := make([]byte, 1)
	for {
		_, err = file.Read(current_symbol)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if current_symbol[0] == '\r' {
			continue
		}
		frequency[current_symbol[0]]++
	}
	forest := make([]models.Forest, len(frequency))
	symbols := make([]byte, len(frequency))
	i := 0
	for symbol, value := range frequency {
		forest[i] = models.Forest{Weight: value, Root: i}
		symbols[i] = symbol
		i++
	}
	return forest, symbols, nil
}

func minMin(forest []models.Forest, forest_size int) (int, int) {
	min1, min2 := math.MaxInt, math.MaxInt
	min1_index, min2_index := -1, -1
	for i, v := range forest {
		if forest_size <= 0 {
			break
		}
		if v.Weight < min1 {
			min2 = min1
			min2_index = min1_index
			min1 = v.Weight
			min1_index = i
		} else if v.Weight < min2 {
			min2 = v.Weight
			min2_index = i
		}
	}
	return min1_index, min2_index
}

func createTree(forest []models.Forest, symbols []byte) []models.Node {
	tree := make([]models.Node, 2*len(forest)-1)
	for i, f := range forest {
		node := models.NewNode()
		node.Symbol = symbols[i]
		tree[f.Root] = *node
	}
	forest_size := len(forest)
	tree_size := forest_size

	for forest_size > 1 {
		min1_index, min2_index := minMin(forest, forest_size)

		tree[tree_size] = models.Node{Left: forest[min1_index].Root, Right: forest[min2_index].Root, Parent: -1}
		tree[forest[min1_index].Root].Parent = tree_size
		tree[forest[min2_index].Root].Parent = tree_size

		forest[min1_index].Root = tree_size
		forest[min2_index].Root = tree_size
		forest[min1_index].Weight += forest[min2_index].Weight

		forest = append(forest[:min2_index], forest[min2_index+1:]...)
		tree_size++
		forest_size--
	}
	return tree
}

func generateCodes(tree []models.Node) map[byte]byte {
	codes := make(map[byte]byte)

	var traverse func(nodeIndex int, currentCode byte)
	traverse = func(nodeIndex int, currentCode byte) {
		node := tree[nodeIndex]

		if node.Left == -1 && node.Right == -1 {
			codes[node.Symbol] = currentCode
			return
		}
		if node.Left != -1 {
			traverse(node.Left, currentCode*2)
		}
		if node.Right != -1 {
			traverse(node.Right, currentCode*2+1)
		}
	}

	if len(tree) > 0 {
		traverse(len(tree)-1, 0)
	}
	return codes
}

func encodeFile(inputFileName, outputFileName string, codes map[byte]byte) error {
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	output := make([]byte, 0)
	currentSymbol := make([]byte, 1)
	var currentByte byte
	currentByteSize := 0
	for {
		_, err = inputFile.Read(currentSymbol)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		code := codes[currentSymbol[0]]
		codeSize := len(strconv.FormatInt(int64(code), 2))
		for {
			if codeSize == 0 {
				break
			}
			if currentByteSize == 8 {
				output = append(output, currentByte)
				currentByteSize = 0
				currentByte = 0
			}
			currentByte = currentByte*2 + code%2
			code /= 2
			currentByteSize++
			codeSize--
		}
		if codeSize == 0 {
			continue
		}
	}
	if currentByteSize > 0 {
		lastByteSize := currentByteSize
		for currentByteSize < 8 {
			currentByteSize++
			currentByte *= 2
		}
		var maskByte byte
		for range 8 {
			maskByte *= 2
			if lastByteSize > 0 {
				lastByteSize--
				maskByte++
			}
		}
		output = append(output, currentByte, maskByte)
	}
	outputFile, err := os.OpenFile(outputFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = outputFile.Write(output)
	if err != nil {
		return err
	}
	var testOutput string
	for _, i := range output {
		testOutput += fmt.Sprintf("%b ", i)
	}
	// fmt.Println(testOutput)
	return nil
}

func encodeTree(outputFileName string, tree []models.Node) error {
	outputFile, err := os.OpenFile(outputFileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	outputFile.Write([]byte(fmt.Sprintf("%d\n", len(tree))))

	for _, node := range tree {
		_, err = outputFile.WriteString(fmt.Sprintf("%d %d %d %d\n", node.Left, node.Right, node.Parent, node.Symbol))
		if err != nil {
			return err
		}
	}
	return nil
}

func Archive(inputFileName, outputFileName string) error {
	forest, symbols, err := createForest(inputFileName)
	if err != nil {
		return err
	}

	tree := createTree(forest, symbols)

	codes := generateCodes(tree)

	err = encodeTree(outputFileName, tree)
	if err != nil {
		return err
	}
	err = encodeFile(inputFileName, outputFileName, codes)
	if err != nil {
		return err
	}

	return nil
}
