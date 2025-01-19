package main

// Queue очередь на основе слайса.
type Queue struct {
	sl []Cell
}

// Enqueue добавление элементов в конец очереди.
func (q *Queue) Enqueue(cells []Cell) {
	q.sl = append(q.sl, cells...)
}

// IsEmpty возвращает true если очередь пуста и false иначе.
func (q *Queue) IsEmpty() bool {
	return len(q.sl) == 0
}

// Dequeue возвращает элемент из начала очереди. При этом выбранный элемент из очереди
// удаляется.
func (q *Queue) Dequeue() (Cell, bool) {
	if q.IsEmpty() {
		return Cell{}, false //nolint:exhaustruct
	}
	first := q.sl[0]
	q.sl = q.sl[1:len(q.sl)]

	return first, true
}
