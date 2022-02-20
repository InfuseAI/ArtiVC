package executor

import (
	"context"
	"runtime"
	"sync"
)

type TaskFunc func(ctx context.Context) error

func ExecuteAll(numCPU int, tasks ...TaskFunc) error {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(len(tasks))

	if numCPU == 0 {
		numCPU = runtime.NumCPU()
	}
	queue := make(chan TaskFunc, numCPU)

	// Spawn the executer
	for i := 0; i < numCPU; i++ {
		go func() {
			for task := range queue {
				if err == nil {
					taskErr := task(ctx)
					if taskErr != nil {
						err = taskErr
						cancel()
					}
				}
				wg.Done()
			}
		}()
	}

	// Add tasks to queue
	for _, task := range tasks {
		queue <- task
	}
	close(queue)

	// wait for all task done
	wg.Wait()
	return err
}
