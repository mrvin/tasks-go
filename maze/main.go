package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type Cell struct {
	i int
	j int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int // Длина и ширина лабиринта.
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		log.Print("can't scan length and width of maze")
		return
	}

	maze := make([][]uint8, n) // Лабиринт - двухмерный массив (матрица).
	for i := range n {
		maze[i] = make([]uint8, m)
		for j := range m {
			if _, err := fmt.Fscan(in, &maze[i][j]); err != nil {
				log.Print("can't scan maze")
				return
			}
		}
	}

	var start, finish Cell // Старт и финиш.
	if _, err := fmt.Fscan(in, &start.i, &start.j); err != nil {
		log.Print("can't scan start cell")
		return
	}
	if _, err := fmt.Fscan(in, &finish.i, &finish.j); err != nil {
		log.Print("can't scan finish cell")
		return
	}

	paths, sums := BFS(maze, start, finish)

	minSum := math.MaxInt
	indexMinSum := -1
	for i, sum := range sums {
		if sum < minSum {
			indexMinSum = i
			minSum = sum
		}
	}
	if minSum == math.MaxInt && indexMinSum == -1 {
		log.Print("path does not exist")
		return
	}
	for _, c := range paths[indexMinSum] {
		fmt.Fprintf(out, "%d %d\n", c.i, c.j)
	}
	fmt.Fprintf(out, ".\n")
	fmt.Fprintf(out, "sum: %d\n", minSum)
}
