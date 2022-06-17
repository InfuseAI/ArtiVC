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

	if numCPU == 0 {
		numCPU = runtime.NumCPU()
	}

	wg := sync.WaitGroup{}
	wg.Add(numCPU)

	queue := make(chan TaskFunc, len(tasks))
	// Add tasks to queue
	for _, task := range tasks {
		queue <- task
	}
	close(queue)

	// Spawn the executer
	for i := 0; i < numCPU; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-queue:
					if ctx.Err() != nil || !ok {
						return
					}
					if e := task(ctx); e != nil {
						err = e
						cancel()
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// wait for all task done
	wg.Wait()
	return err
}
