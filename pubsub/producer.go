package pubsub

type Producer interface {
	Produce(task interface{})
}
