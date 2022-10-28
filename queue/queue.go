package queue

type Queue struct {
	Length int64
	items  chan interface{}
}

func (q *Queue) Push(value interface{}) {
	go func() {
		q.items <- value
	}()
}

func (q *Queue) Pop() (result interface{}) {
	result = <-q.items
	return result
}

func New(len int64) *Queue {
	return &Queue{
		Length: len,
		items:  make(chan interface{}),
	}
}
