package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	inputFileName  = "input.txt"
	outputFileName = "output.txt"
)

type Forest struct {
	Weight int
	Root   int
}

type Node struct {
	Left   int
	Right  int
	Parent int
	Symbol byte
}

func NewNode() *Node {
	return &Node{Left: -1, Right: -1, Parent: -1}
}

func CreateForest(inputFileName string) ([]Forest, []byte, error) {
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
	forest := make([]Forest, len(frequency))
	symbols := make([]byte, len(frequency))
	i := 0
	for symbol, value := range frequency {
		forest[i] = Forest{Weight: value, Root: i}
		symbols[i] = symbol
		i++
	}
	return forest, symbols, nil
}

func MinMin(forest []Forest, forest_size int) (int, int) {
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

func CreateTree(forest []Forest, symbols []byte) []Node {
	tree := make([]Node, 2*len(forest)-1)
	for i, f := range forest {
		node := NewNode()
		node.Symbol = symbols[i]
		tree[f.Root] = *node
	}
	forest_size := len(forest)
	tree_size := forest_size

	for forest_size > 1 {
		min1_index, min2_index := MinMin(forest, forest_size)

		tree[tree_size] = Node{Left: forest[min1_index].Root, Right: forest[min2_index].Root, Parent: -1}
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

func GenerateCodes(tree []Node) map[byte]byte {
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

func EncodeFromString(input string, codes map[byte]byte) string {
	var encoded string
	for _, c := range []byte(input) {
		encoded += fmt.Sprintf("%b", codes[c])
	}
	return encoded
}

func EncodeFromFile(inputFileName, outputFileName string, codes map[byte]byte) error {
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
	fmt.Println(testOutput)
	return nil
}

func DecodeFromString(encoded string, tree []Node) string {
	var decoded string
	currentNode := len(tree) - 1

	for _, bit := range encoded {
		if bit == '0' {
			currentNode = tree[currentNode].Left
		} else {
			currentNode = tree[currentNode].Right
		}

		if tree[currentNode].Left == -1 && tree[currentNode].Right == -1 {
			decoded += string(tree[currentNode].Symbol)
			currentNode = len(tree) - 1
		}
	}
	return decoded
}

func EncodeTree(tree []Node, outputFileName string) error {
	outputFile, err := os.OpenFile(outputFileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	outputFile.Write([]byte(fmt.Sprintf("%d\n", len(tree))))

	for i, node := range tree {
		_, err = outputFile.WriteString(fmt.Sprintf("%d %d %d %d %d\n", i, node.Left, node.Right, node.Parent, node.Symbol))
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	forest, symbols, err := CreateForest(inputFileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	tree := CreateTree(forest, symbols)

	codes := GenerateCodes(tree)

	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	input := make([]byte, 1024)
	n, err := file.Read(input)
	if err != nil {
		log.Fatal(err)
		return
	}

	for i := 0; i < n; i++ {
		if input[i] == '\r' {
			input = append(input[:i], input[i+1:]...)
			n--
			i--
		}
	}

	encoded := EncodeFromString(string(input[:n]), codes)
	decoded := DecodeFromString(encoded, tree)

	// fmt.Printf("encoded: %v\n", encoded)
	fmt.Printf("decoded: %s\n", decoded)
	fmt.Printf("input:   %s\n", string(input[:n]))

	// for i := range decoded {
	// 	if decoded[i] != input[i] {
	// 		fmt.Printf("i: %d, decoded: %d, input: %d\n", i, byte(decoded[i]), byte(input[i]))
	// 	}
	// }

	err = EncodeTree(tree, outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	err = EncodeFromFile(inputFileName, outputFileName, codes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encoded)
}
