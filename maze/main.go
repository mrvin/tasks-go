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

const minMazeLenWidth = 1
const maxMazeCell = 9

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	var n, m int // Длина и ширина лабиринта.
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		log.Fatal("can't scan length and width of maze")
	}
	if n < minMazeLenWidth || m < minMazeLenWidth {
		log.Fatalf("length or width < %d", minMazeLenWidth)
	}

	maze := make([][]uint8, n) // Лабиринт - двухмерный массив (матрица).
	for i := range n {
		maze[i] = make([]uint8, m)
		for j := range m {
			if _, err := fmt.Fscan(in, &maze[i][j]); err != nil {
				log.Fatal("can't scan maze")
			}
			if maze[i][j] > maxMazeCell {
				log.Fatalf("maze cell  > %d", maxMazeCell)
			}
		}
	}

	var start, finish Cell // Старт и финиш.
	if _, err := fmt.Fscan(in, &start.i, &start.j); err != nil {
		log.Fatal("can't scan start cell")
	}
	if start.i >= n && start.j >= m {
		log.Fatal("start outside maze")
	}
	if _, err := fmt.Fscan(in, &finish.i, &finish.j); err != nil {
		log.Fatal("can't scan finish cell")
	}
	if finish.i >= n && finish.j >= m {
		log.Fatal("finish outside maze")
	}

	path, pathLen := LeeWithWeight(maze, start, finish)

	for _, c := range path {
		fmt.Fprintf(out, "%d %d\n", c.i, c.j)
	}
	fmt.Fprintf(out, ".\n")
	fmt.Fprintf(out, "len: %d\n", pathLen)
	out.Flush()
}
