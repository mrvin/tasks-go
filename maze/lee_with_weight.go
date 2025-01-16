package main

import (
	"math"
)

// visit возвращает слайс ячеек в которые доступен переход из ячейки stand.
func visit(maze [][]uint8, mazeMinSum [][]int, stand Cell, minSum int) []Cell {
	i := stand.i
	j := stand.j
	cells := make([]Cell, 0)

	// Вверх
	if i-1 >= 0 && maze[i-1][j] != 0 {
		newMinSum := mazeMinSum[i][j] + int(maze[i-1][j])
		if (mazeMinSum[i-1][j] == 0 || newMinSum < mazeMinSum[i-1][j]) && newMinSum < minSum {
			mazeMinSum[i-1][j] = newMinSum
			cells = append(cells, Cell{i - 1, j})
		}
	}
	// Вправо
	if j+1 < len(maze[i]) && maze[i][j+1] != 0 {
		newMinSum := mazeMinSum[i][j] + int(maze[i][j+1])
		if (mazeMinSum[i][j+1] == 0 || newMinSum < mazeMinSum[i][j+1]) && newMinSum < minSum {
			mazeMinSum[i][j+1] = newMinSum
			cells = append(cells, Cell{i, j + 1})
		}
	}
	// Влево
	if j-1 >= 0 && maze[i][j-1] != 0 {
		newMinSum := mazeMinSum[i][j] + int(maze[i][j-1])
		if (mazeMinSum[i][j-1] == 0 || newMinSum < mazeMinSum[i][j-1]) && newMinSum < minSum {
			mazeMinSum[i][j-1] = newMinSum
			cells = append(cells, Cell{i, j - 1})
		}
	}
	// Вниз
	if i+1 < len(maze) && maze[i+1][j] != 0 {
		newMinSum := mazeMinSum[i][j] + int(maze[i+1][j])
		if (mazeMinSum[i+1][j] == 0 || newMinSum < mazeMinSum[i+1][j]) && newMinSum < minSum {
			mazeMinSum[i+1][j] = newMinSum
			cells = append(cells, Cell{i + 1, j})
		}
	}

	return cells
}

// nextMinSumCell возвращает соседнюю ячейку с минимальным растоянием до финиша и ok
// равный true, если такой ячейки нет, то возвращает ячейку Cell{i:0, j:0} и ok равный
// false.
func nextMinSumCell(mazeMinSum [][]int, stand Cell) (Cell, bool) {
	i := stand.i
	j := stand.j

	cells := [4]Cell{
		{i - 1, j}, // Вверх
		{i, j + 1}, // Вправо
		{i, j - 1}, // Влево
		{i + 1, j}, // Вниз
	}
	sums := [4]int{math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt}

	if i-1 >= 0 && mazeMinSum[i-1][j] != 0 {
		sums[0] = mazeMinSum[i-1][j]
	}

	if j+1 < len(mazeMinSum[i]) && mazeMinSum[i][j+1] != 0 {
		sums[1] = mazeMinSum[i][j+1]
	}

	if j-1 >= 0 && mazeMinSum[i][j-1] != 0 {
		sums[2] = mazeMinSum[i][j-1]
	}

	if i+1 < len(mazeMinSum) && mazeMinSum[i+1][j] != 0 {
		sums[3] = mazeMinSum[i+1][j]
	}

	minSum := math.MaxInt
	minIndex := 0
	for i, sum := range sums {
		if sum < minSum {
			minSum = sum
			minIndex = i
		}
	}
	if minSum == math.MaxInt {
		return Cell{}, false //nolint:exhaustruct
	}

	return cells[minIndex], true
}

// LeeWithWeight поиск самого быстрого пути в лабиринте при помощи алгоритма Ли
// (волновой алгоритм) с учетом длины пути (веса) на основе поиска в ширину(BFS).
// Возвращает самый быстрый путь и его длину, если таких путей несколько, то возвращает
// случайный. Если пути нет, то возвращает слайс нулевой длины и длину равную 0.
func LeeWithWeight(maze [][]uint8, start, finish Cell) ([]Cell, int) {
	// Инициализация
	mazeMinSum := make([][]int, len(maze))
	for i := range mazeMinSum {
		mazeMinSum[i] = make([]int, len(maze[i]))
	}
	resultMinSum := math.MaxInt

	// Распространение волны
	var queue Queue
	mazeMinSum[finish.i][finish.j] = int(maze[finish.i][finish.j])
	queue.Enqueue([]Cell{finish})
	for !queue.IsEmpty() {
		stand, _ := queue.Dequeue()
		if stand == start {
			resultMinSum = mazeMinSum[start.i][start.j]
			continue
		}
		cells := visit(maze, mazeMinSum, stand, resultMinSum)
		queue.Enqueue(cells)
	}

	// Восстановление пути
	var resultPaths []Cell
	stand := start
	for {
		resultPaths = append(resultPaths, stand)
		if stand == finish {
			break
		}
		var ok bool
		stand, ok = nextMinSumCell(mazeMinSum, stand)
		if !ok {
			return []Cell{}, 0
		}
	}

	return resultPaths, resultMinSum
}
