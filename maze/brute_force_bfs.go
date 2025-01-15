package main

import (
	"math"
	"slices"
)

type CellPathSum struct {
	path []Cell
	sum  int
}

type Queue struct {
	sl []CellPathSum
}

func (q *Queue) Enqueue(cells []CellPathSum) {
	q.sl = append(q.sl, cells...)
}

func (q *Queue) IsEmpty() bool {
	return len(q.sl) == 0
}

func (q *Queue) Dequeue() (CellPathSum, bool) {
	if q.IsEmpty() {
		return CellPathSum{}, false //nolint:exhaustruct
	}
	first := q.sl[0]
	q.sl = q.sl[1:len(q.sl)]

	return first, true
}

func visit(maze [][]uint8, stand CellPathSum, minSum int) []CellPathSum {
	i := stand.path[len(stand.path)-1].i
	j := stand.path[len(stand.path)-1].j
	cells := make([]CellPathSum, 0)

	// Вверх
	if i-1 >= 0 && maze[i-1][j] > 0 && !slices.Contains(stand.path, Cell{i - 1, j}) {
		newPath := append(stand.path[0:len(stand.path):len(stand.path)], Cell{i - 1, j})
		newSum := stand.sum + int(maze[i-1][j])
		if newSum < minSum {
			cells = append(cells, CellPathSum{newPath, newSum})
		}
	}
	// Вправо
	if j+1 < len(maze[i]) && maze[i][j+1] > 0 && !slices.Contains(stand.path, Cell{i, j + 1}) {
		newPath := append(stand.path[0:len(stand.path):len(stand.path)], Cell{i, j + 1})
		newSum := stand.sum + int(maze[i][j+1])
		if newSum < minSum {
			cells = append(cells, CellPathSum{newPath, newSum})
		}
	}
	// Влево
	if j-1 >= 0 && maze[i][j-1] > 0 && !slices.Contains(stand.path, Cell{i, j - 1}) {
		newPath := append(stand.path[0:len(stand.path):len(stand.path)], Cell{i, j - 1})
		newSum := stand.sum + int(maze[i][j-1])
		if newSum < minSum {
			cells = append(cells, CellPathSum{newPath, newSum})
		}
	}
	// Вниз
	if i+1 < len(maze) && maze[i+1][j] > 0 && !slices.Contains(stand.path, Cell{i + 1, j}) {
		newPath := append(stand.path[0:len(stand.path):len(stand.path)], Cell{i + 1, j})
		newSum := stand.sum + int(maze[i+1][j])
		if newSum < minSum {
			cells = append(cells, CellPathSum{newPath, newSum})
		}
	}

	return cells
}

// Перебор всех возможных путей на основе поиска в ширину. Если уже найден путь к
// финишу, то все пути с большей или равной длиной пути будут отброшены.
func BruteForceBFS(maze [][]uint8, start, finish Cell) ([]Cell, int) {
	var resultPaths []Cell
	resultSum := math.MaxInt

	var queue Queue
	queue.Enqueue([]CellPathSum{{[]Cell{start}, int(maze[start.i][start.j])}})
	for !queue.IsEmpty() {
		stand, _ := queue.Dequeue()
		if stand.path[len(stand.path)-1] == finish {
			resultPaths = stand.path
			resultSum = stand.sum
			continue
		}
		cells := visit(maze, stand, resultSum)
		queue.Enqueue(cells)
	}

	return resultPaths, resultSum
}
