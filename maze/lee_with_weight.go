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
	mazeMinPathLen := make([][]int, len(maze))
	for i := range mazeMinPathLen {
		mazeMinPathLen[i] = make([]int, len(maze[i]))
	}
	resultMinPathLen := math.MaxInt

	// Распространение волны
	var queue Queue
	mazeMinPathLen[finish.i][finish.j] = int(maze[finish.i][finish.j])
	queue.Enqueue([]Cell{finish})
	for !queue.IsEmpty() {
		stand, _ := queue.Dequeue()
		if stand == start {
			resultMinPathLen = mazeMinPathLen[start.i][start.j]
			continue
		}
		cells := visit(maze, mazeMinPathLen, stand, resultMinPathLen)
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
		stand, ok = nextMinPathLenCell(mazeMinPathLen, stand)
		if !ok {
			return []Cell{}, 0 // Путь не найден.
		}
	}

	return resultPath, resultMinPathLen
}

// visit возвращает слайс клеток в которые доступен переход из клетки stand.
func visit(maze [][]uint8, mazeMinPathLen [][]int, stand Cell, minPathLen int) []Cell {
	i := stand.i
	j := stand.j
	cells := make([]Cell, 0)

	// Вверх
	if i-1 >= 0 && maze[i-1][j] != 0 {
		newMinPathLen := mazeMinPathLen[i][j] + int(maze[i-1][j])
		if (mazeMinPathLen[i-1][j] == 0 || newMinPathLen < mazeMinPathLen[i-1][j]) && newMinPathLen < minPathLen {
			mazeMinPathLen[i-1][j] = newMinPathLen
			cells = append(cells, Cell{i - 1, j})
		}
	}
	// Вправо
	if j+1 < len(maze[i]) && maze[i][j+1] != 0 {
		newMinPathLen := mazeMinPathLen[i][j] + int(maze[i][j+1])
		if (mazeMinPathLen[i][j+1] == 0 || newMinPathLen < mazeMinPathLen[i][j+1]) && newMinPathLen < minPathLen {
			mazeMinPathLen[i][j+1] = newMinPathLen
			cells = append(cells, Cell{i, j + 1})
		}
	}
	// Влево
	if j-1 >= 0 && maze[i][j-1] != 0 {
		newMinPathLen := mazeMinPathLen[i][j] + int(maze[i][j-1])
		if (mazeMinPathLen[i][j-1] == 0 || newMinPathLen < mazeMinPathLen[i][j-1]) && newMinPathLen < minPathLen {
			mazeMinPathLen[i][j-1] = newMinPathLen
			cells = append(cells, Cell{i, j - 1})
		}
	}
	// Вниз
	if i+1 < len(maze) && maze[i+1][j] != 0 {
		newMinPathLen := mazeMinPathLen[i][j] + int(maze[i+1][j])
		if (mazeMinPathLen[i+1][j] == 0 || newMinPathLen < mazeMinPathLen[i+1][j]) && newMinPathLen < minPathLen {
			mazeMinPathLen[i+1][j] = newMinPathLen
			cells = append(cells, Cell{i + 1, j})
		}
	}

	return cells
}

// nextMinSumCell возвращает соседнюю c stand клетку с минимальным растоянием до
// финиша и ok равный true, если такой клетки нет, то возвращает клетку Cell{i:0, j:0}
// и ok равный false.
func nextMinPathLenCell(mazeMinPathLen [][]int, stand Cell) (Cell, bool) {
	i := stand.i
	j := stand.j

	cells := [4]Cell{
		{i - 1, j}, // Вверх
		{i, j + 1}, // Вправо
		{i, j - 1}, // Влево
		{i + 1, j}, // Вниз
	}
	lens := [4]int{math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt}

	if i-1 >= 0 && mazeMinPathLen[i-1][j] != 0 {
		lens[0] = mazeMinPathLen[i-1][j]
	}

	if j+1 < len(mazeMinPathLen[i]) && mazeMinPathLen[i][j+1] != 0 {
		lens[1] = mazeMinPathLen[i][j+1]
	}

	if j-1 >= 0 && mazeMinPathLen[i][j-1] != 0 {
		lens[2] = mazeMinPathLen[i][j-1]
	}

	if i+1 < len(mazeMinPathLen) && mazeMinPathLen[i+1][j] != 0 {
		lens[3] = mazeMinPathLen[i+1][j]
	}

	minLen := math.MaxInt
	index := 0
	for i, l := range lens {
		if l < minLen {
			index = i
			minLen = l
		}
	}
	if minLen == math.MaxInt {
		return Cell{}, false //nolint:exhaustruct
	}

	return cells[index], true
}
