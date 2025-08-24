package events

import (
	"log"
	"mini-alt/events/types"
	"sync"
)

type WorkerPool struct {
	Jobs chan WebhookJob
	wg   sync.WaitGroup
}

type WebhookJob struct {
	Event *types.Event
	Url   string
	Token string
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		Jobs: make(chan WebhookJob, 1000),
	}

	for i := 0; i < numWorkers; i++ {
		go pool.worker(i)
	}

	return pool
}

func (p *WorkerPool) worker(id int) {
	for job := range p.Jobs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Worker %d] Recovered from panic: %v", id, r)
					p.wg.Done()
				}
			}()

			err := sendWebhook(job.Url, job.Token, job.Event)
			if err != nil {
				log.Printf("[Worker %d] Failed to send webhook event to webhook server: %v", id, err)
			}

			p.wg.Done()
		}()
	}
}

func (p *WorkerPool) Submit(e WebhookJob) {
	if e.Event == nil {
		log.Println("WARNING: nil event submitted")
	}
	p.wg.Add(1)
	p.Jobs <- e
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}
