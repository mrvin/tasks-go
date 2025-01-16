package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Cell struct {
	i int
	j int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	var n, m int // Длина и ширина лабиринта.
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		log.Fatal("can't scan length and width of maze")
	}

	maze := make([][]uint8, n) // Лабиринт - двухмерный массив (матрица).
	for i := range n {
		maze[i] = make([]uint8, m)
		for j := range m {
			if _, err := fmt.Fscan(in, &maze[i][j]); err != nil {
				log.Fatal("can't scan maze")
			}
		}
	}

	var start, finish Cell // Старт и финиш.
	if _, err := fmt.Fscan(in, &start.i, &start.j); err != nil {
		log.Fatal("can't scan start cell")
	}
	if _, err := fmt.Fscan(in, &finish.i, &finish.j); err != nil {
		log.Fatal("can't scan finish cell")
	}

	path, sum := LeeWithWeight(maze, start, finish)

	for _, c := range path {
		fmt.Fprintf(out, "%d %d\n", c.i, c.j)
	}
	fmt.Fprintf(out, ".\n")
	fmt.Fprintf(out, "sum: %d\n", sum)
	out.Flush()
}
