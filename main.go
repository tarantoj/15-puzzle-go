package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type node struct {
	tiles  []int
	g      int // cost to get here
	parent *node
	a      action
	index  int
}

func (n node) length() int {
	return len(n.tiles)
}

func (n node) width() int {
	return 4
}

func (n node) heuristic() int {
	dist := 0

	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x

	}

	for i, v := range n.tiles {
		dist += abs(i%n.width() - v%n.width())
		dist += abs(i/n.width() - v/n.width())
	}
	return dist
}

func (n node) zero() int {
	var zero int
	for i, v := range n.tiles {
		if v == 0 {
			zero = i
			break
		}
	}
	return zero
}

func (n node) f() int {
	return n.g + n.heuristic()
}

type action uint8

const (
	up action = 1 << iota
	right
	down
	left
)

func (n node) permutations() []node {

	var nodes []node

	create := func(n node) node {
		var new node
		tmp := make([]int, n.length())
		copy(tmp, n.tiles)
		new.tiles = tmp
		new.g = n.g + 1
		new.parent = &n
	}

	if n.zero() > 3 {
		new := create(n)
		new.tiles[n.zero()] = new.tiles[n.zero()-4]
		new.tiles[n.zero()-4] = 0
		new.a = up
		nodes = append(nodes, new)
	}

	if (n.zero()+1)%4 != 0 {
		new := create(n)
		new.tiles[n.zero()] = new.tiles[n.zero()+1]
		new.tiles[n.zero()+1] = 0
		new.a = right
		new.g = n.g + 1
		nodes = append(nodes, new)
	}

	if n.zero() < 12 {
		new := create(n)
		new.tiles[n.zero()] = new.tiles[n.zero()+4]
		new.tiles[n.zero()+4] = 0
		new.a = down
		new.g = n.g + 1
		nodes = append(nodes, new)
	}

	if n.zero()%4 != 0 {
		new := create(n)
		new.tiles[n.zero()] = new.tiles[n.zero()-1]
		new.tiles[n.zero()-1] = 0
		new.a = left
		new.g = n.g + 1
		nodes = append(nodes, new)
	}

	return nodes
}

func printPuzzle(n *node) {
	for i := range n.tiles {
		fmt.Printf("%2d ", n.tiles[i])
		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to lowest cost, so we use less than.
	return pq[i].f() > pq[j].f()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push an item to q
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*node)
	item.index = n
	*pq = append(*pq, item)
}

// Pop item from q
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func astar(n node) {

}

func main() {
	var (
		next      int
		puzzle    []int
		inputfile string
		root      node
	)

	flag.StringVar(&inputfile, "i", "", "Input file")
	flag.Parse()

	file, err := os.Open(inputfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for {
		n, err := fmt.Fscanf(file, "%d", &next)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		if n == 0 {
			break
		}

		puzzle = append(puzzle, next)
	}

	if len(puzzle)%4 != 0 {
		panic("Puzzle must be square")
	}

	root.tiles = puzzle
	root.g = 0
	fmt.Printf("%+v\n", root.permutations())

	return
}
