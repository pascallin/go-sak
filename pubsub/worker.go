package pubsub

type Worker struct {
	jobs   chan *models.Task
	Action func() error
	WorkerInfo
}

type WorkerInfo struct {
	FailedJobs    map[string]*models.Task `json:"failedJobs"`
	CompletedJobs map[string]*models.Task `json:"completedJobs"`
	RunningJobs   map[string]*models.Task `json:"runningJobs"`
}

func NewWorker() *Worker {
	worker := &Worker{
		jobs:   make(chan *models.Task),
		Action: func() error { fmt.Println("empty action...") },
		WorkerInfo: WorkerInfo{
			FailedJobs:    make(map[string]*models.Task),
			CompletedJobs: make(map[string]*models.Task),
			RunningJobs:   make(map[string]*models.Task),
		},
	}

	return worker
}

func (w *Worker) GetInfo() WorkerInfo {
	return w.WorkerInfo
}

func (w *Worker) Close() {
	close(w.jobs)
}

func (w *Worker) Produce(task *models.Task) {
	w.jobs <- task
}

func (w *Worker) Consume() {
	fmt.Println("worker consuming...")
	for {
		select {
		case t := <-w.jobs:
			fmt.Println("received job", t)
			w.transcode(t)
		}
	}
}
