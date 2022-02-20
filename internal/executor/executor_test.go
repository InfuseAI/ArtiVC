package executor

import (
	"context"
	"errors"
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHappyPath(t *testing.T) {
	var str1, str2 *string

	task1 := func(ctx context.Context) error {
		str := "foo"
		str1 = &str
		return nil
	}

	task2 := func(ctx context.Context) error {
		str := "bar"
		str2 = &str
		return nil
	}

	err := ExecuteAll(runtime.NumCPU(), task1, task2)
	assert.Empty(t, err)
	assert.Equal(t, "foo", *str1)
	assert.Equal(t, "bar", *str2)
}

func TestFailedPath(t *testing.T) {
	ErrFoo := errors.New("foo")

	taskOk := func(ctx context.Context) error {
		return nil
	}

	taskErr := func(ctx context.Context) error {
		return ErrFoo
	}

	err := ExecuteAll(runtime.NumCPU(), taskOk, taskErr)
	assert.Equal(t, ErrFoo, err)

	err = ExecuteAll(runtime.NumCPU(), taskErr, taskOk)
	assert.Equal(t, ErrFoo, err)
}

func TestConcurrent(t *testing.T) {
	tasks := []TaskFunc{}
	var counter int32

	for i := 0; i < 100; i++ {
		f := func(ctx context.Context) error {
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			atomic.AddInt32(&counter, 1)
			return nil
		}
		tasks = append(tasks, f)

	}

	err := ExecuteAll(50, tasks...)
	assert.Empty(t, err)
	assert.Equal(t, int32(100), counter)
}

func TestContext(t *testing.T) {
	ErrFoo := errors.New("foo")

	taskForever := func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}

	taskErr := func(ctx context.Context) error {
		return ErrFoo
	}

	err := ExecuteAll(3, taskForever, taskErr)
	assert.Equal(t, ErrFoo, err)
}
