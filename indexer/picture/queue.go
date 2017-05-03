package picture

import (
	"container/list"
	"sync"
)

type Queue struct {
	items *list.List
	mutex sync.Mutex
}

func NewQueue() *Queue {
	q := new(Queue)
	q.items = list.New()
	return q
}

func (q *Queue) PushBack(p *Index) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.items.PushBack(p)
}

func (q *Queue) Pop() (int, *Index) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	len := q.items.Len()
	elem := q.items.Front()
	if elem != nil {
		idx := q.items.Front().Value.(*Index)
		q.items.Remove(q.items.Front())
		return len, idx
	}
	return len, nil
}

func (q *Queue) Length() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.items.Len()
}

func (q *Queue) Keys() []string {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	keys := make([]string, q.items.Len())
	elem := q.items.Front()
	i := 0
	for elem != nil {
		idx := elem.Value.(*Index)
		keys[i] = idx.MD5
		elem = elem.Next()
		i++
	}
	return keys
}

func (q *Queue) Items() []*Index {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	items := make([]*Index, q.items.Len())
	elem := q.items.Front()
	i := 0
	for elem != nil {
		idx := elem.Value.(*Index)
		items[i] = idx
		elem = elem.Next()
		i++
	}
	return items
}
