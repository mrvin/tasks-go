package main

import (
	"math"
)

// LeeWithWeight поиск самого быстрого пути в лабиринте при помощи алгоритма Ли
// (волновой алгоритм) с учетом длины пути (веса) на основе поиска в ширину(BFS).
// Возвращает самый быстрый путь и его длину, если таких путей несколько, то
// возвращает случайный. Если пути нет, то возвращает слайс нулевой длины и длину
// равную 0.
func LeeWithWeight(maze [][]uint8, start, finish Cell) ([]Cell, int) {
	// Инициализация
	mazeExitMap := make([][]int, len(maze))
	for i := range mazeExitMap {
		mazeExitMap[i] = make([]int, len(maze[i]))
	}
	resultMinPathLen := math.MaxInt

	// Распространение волны
	var queue Queue
	mazeExitMap[finish.i][finish.j] = int(maze[finish.i][finish.j])
	queue.Enqueue([]Cell{finish})
	for !queue.IsEmpty() {
		stand, _ := queue.Dequeue()
		if stand == start {
			resultMinPathLen = mazeExitMap[start.i][start.j]
			continue
		}
		cells := visitNeighbors(maze, mazeExitMap, stand, resultMinPathLen)
		queue.Enqueue(cells)
	}

	// Восстановление пути
	var resultPath []Cell
	stand := start
	for {
		resultPath = append(resultPath, stand)
		if stand == finish {
			break // Путь найден.
		}
		var ok bool
		stand, ok = nextMinPathLenCell(mazeExitMap, stand)
		if !ok {
			return []Cell{}, 0 // Путь не найден.
		}
	}

	return resultPath, resultMinPathLen
}

var directions = [4][2]int{
	{-1, 0}, // Вверх
	{0, 1},  // Вправо
	{0, -1}, // Влево
	{1, 0},  // Вниз
}

// visitNeighbors возвращает слайс клеток в которые доступен переход из клетки stand.
func visitNeighbors(maze [][]uint8, mazeExitMap [][]int, stand Cell, minPathLen int) []Cell {
	i := stand.i
	j := stand.j
	cells := make([]Cell, 0)

	for _, d := range directions {
		di := i + d[0]
		dj := j + d[1]
		if 0 <= di && di < len(mazeExitMap) && 0 <= dj && dj < len(mazeExitMap[i]) && maze[di][dj] != 0 {
			newMinPathLen := mazeExitMap[i][j] + int(maze[di][dj])
			if (mazeExitMap[di][dj] == 0 || newMinPathLen < mazeExitMap[di][dj]) && newMinPathLen < minPathLen {
				mazeExitMap[di][dj] = newMinPathLen
				cells = append(cells, Cell{di, dj})
			}
		}
	}

	return cells
}

// nextMinPathLenCell возвращает соседнюю c stand клетку с минимальным растоянием до
// финиша и ok равный true, если такой клетки нет, то возвращает клетку Cell{i:0, j:0}
// и ok равный false.
func nextMinPathLenCell(mazeExitMap [][]int, stand Cell) (Cell, bool) {
	i := stand.i
	j := stand.j

	minI, minJ := 0, 0
	minLen := math.MaxInt
	for _, d := range directions {
		di := i + d[0]
		dj := j + d[1]
		if 0 <= di && di < len(mazeExitMap) && 0 <= dj && dj < len(mazeExitMap[i]) {
			if mazeExitMap[di][dj] != 0 && mazeExitMap[di][dj] < minLen {
				minI, minJ = di, dj
				minLen = mazeExitMap[di][dj]
			}
		}
	}

	if minLen == math.MaxInt {
		return Cell{}, false //nolint:exhaustruct
	}

	return Cell{minI, minJ}, true
}
