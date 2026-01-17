package workers

import (
	"log"
	"sync"
	"time"
)

type OrderJob struct {
	OrderID     uint
	OrderNumber string
	UserEmail   string
}

// StartOrderWorkers starts a pool of goroutines
// bufferSize: How many jobs can wait before the handler blocks
// workerCount: How many goroutines work simultaneously
func StartOrderWorkers(bufferSize int, workerCount int, wg *sync.WaitGroup) chan OrderJob {
	jobChannel := make(chan OrderJob, bufferSize)

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, jobChannel, wg)
	}

	log.Printf("Started %d background order workers", workerCount)
	return jobChannel
}

func worker(id int, jobs <-chan OrderJob, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		log.Printf("[Worker %d] Processing order %s...", id, job.OrderNumber)
		processEmail(job)
		processInventoryUpdate(job)
		log.Printf("[Worker %d] Done.", id)
	}
	log.Printf("[Worker %d] Stopping...", id)
}

func processEmail(job OrderJob) {
	// Here would be the real SMTP code
	time.Sleep(2 * time.Second) // Simulate network latency
	log.Printf("Email sent to %s", job.UserEmail)
}

func processInventoryUpdate(job OrderJob) {
	// Here would be a call to the ERP system
	time.Sleep(1 * time.Second)
	log.Printf("Warehouse notified for %s", job.OrderNumber)
}
