package main

type Queue struct {
	sl []Cell
}

func (q *Queue) Enqueue(cells []Cell) {
	q.sl = append(q.sl, cells...)
}

func (q *Queue) IsEmpty() bool {
	return len(q.sl) == 0
}

func (q *Queue) Dequeue() (Cell, bool) {
	if q.IsEmpty() {
		return Cell{}, false //nolint:exhaustruct
	}
	first := q.sl[0]
	q.sl = q.sl[1:len(q.sl)]

	return first, true
}
