package main

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

func (q *Queue) Size() int {
	return len(q.sl)
}

func (q *Queue) Dequeue() (CellPathSum, bool) {
	if q.IsEmpty() {
		return CellPathSum{}, false //nolint:exhaustruct
	}
	first := q.sl[0]
	q.sl = q.sl[1:len(q.sl)]

	return first, true
}

func visit(maze [][]uint8, stand CellPathSum) []CellPathSum {
	i := stand.path[len(stand.path)-1].i
	j := stand.path[len(stand.path)-1].j
	cells := make([]CellPathSum, 0)

	if i-1 >= 0 && maze[i-1][j] > 0 && maze[i-1][j] < 10 {
		newPath := append(stand.path[0:len(stand.path):len(stand.path)], Cell{i - 1, j})
		cells = append(cells, CellPathSum{
			newPath,
			stand.sum + int(maze[i-1][j])},
		)
		maze[i-1][j] = 10 // mark as visited
	}
	if j-1 >= 0 && maze[i][j-1] > 0 && maze[i][j-1] < 10 {
		newPath := append(stand.path[0:len(stand.path):len(stand.path)], Cell{i, j - 1})
		cells = append(cells, CellPathSum{
			newPath,
			stand.sum + int(maze[i][j-1])},
		)
		maze[i][j-1] = 10 // mark as visited
	}
	if j+1 < len(maze[i]) && maze[i][j+1] > 0 && maze[i][j+1] < 10 {
		newPath := append(stand.path[0:len(stand.path):len(stand.path)], Cell{i, j + 1})
		cells = append(cells, CellPathSum{
			newPath,
			stand.sum + int(maze[i][j+1])},
		)
		maze[i][j+1] = 10 // mark as visited
	}
	if i+1 < len(maze) && maze[i+1][j] > 0 && maze[i+1][j] < 10 {
		newPath := append(stand.path[0:len(stand.path):len(stand.path)], Cell{i + 1, j})
		cells = append(cells, CellPathSum{
			newPath,
			stand.sum + int(maze[i+1][j])},
		)
		maze[i+1][j] = 10 // mark as visited
	}

	return cells
}

// Поиск в ширину.
func BFS(maze [][]uint8, start, finish Cell) ([][]Cell, []int) {
	var resultPaths [][]Cell
	var resultSum []int

	var queue Queue
	queue.Enqueue([]CellPathSum{{[]Cell{start}, int(maze[start.i][start.j])}})
	maze[start.i][start.j] = 10 // mark as visited
	for !queue.IsEmpty() {
		queueLen := queue.Size()
		for range queueLen {
			stand, _ := queue.Dequeue()
			if stand.path[len(stand.path)-1] == finish {
				resultPaths = append(resultPaths, stand.path)
				resultSum = append(resultSum, stand.sum)
				continue
			}
			cells := visit(maze, stand)
			queue.Enqueue(cells)
		}
	}

	return resultPaths, resultSum
}
