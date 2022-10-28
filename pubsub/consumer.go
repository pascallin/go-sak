package pubsub

type Consumer interface {
	Consume()
}

type ConsumerWorker interface {
}
