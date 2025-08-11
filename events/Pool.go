package events

var Pool *WorkerPool

func InitPool(workers int) {
	Pool = NewWorkerPool(workers)
}
