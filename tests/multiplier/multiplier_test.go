package multiplier

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jit-brains/go/multiplier"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func Test_Multiplier(test *testing.T) {
	test.Run(
		"Create multiplier with 10 worker count. Give them 5 messages and stop multiplier. Then only 5 workers make action",
		func(test *testing.T) {
			var actionCount atomic.Uint64
			taskChannel := make(chan struct{}, 5)
			worker := func() {
				<-taskChannel
				actionCount.Add(1)
			}

			multiplier := multiplier.New(worker).
				WithWorkerCount(5).
				Run()

			taskChannel <- struct{}{}
			taskChannel <- struct{}{}
			taskChannel <- struct{}{}
			taskChannel <- struct{}{}
			taskChannel <- struct{}{}
			multiplier.Stop()
			taskChannel <- struct{}{}
			taskChannel <- struct{}{}

			assert.Equal(test, uint64(5), actionCount.Load())
			goleak.VerifyNone(test)
		},
	)

	test.Run(
		"Create multiplier with default workers count. Give 3 messages twice. At first it's become to panic. At second - ok, cause workers count maintain automatically.",
		func(test *testing.T) {
			var actionCount atomic.Uint64
			taskChannel := make(chan struct{}, 5)
			isNeedPanic := true
			lock := sync.Mutex{}
			worker := func() {
				_, isOpen := <-taskChannel
				if !isOpen {
					return
				}

				lock.Lock()
				isNeedPanic := isNeedPanic
				lock.Unlock()
				if isNeedPanic {
					panic("test panic")
				}

				actionCount.Add(1)
			}

			multiplier := multiplier.New(worker).
				WithMultiplyInterval(10 * time.Millisecond).
				Run()

			taskChannel <- struct{}{}
			taskChannel <- struct{}{}
			taskChannel <- struct{}{}
			time.Sleep(25 * time.Millisecond)

			lock.Lock()
			isNeedPanic = false
			lock.Unlock()

			taskChannel <- struct{}{}
			taskChannel <- struct{}{}
			taskChannel <- struct{}{}
			time.Sleep(35 * time.Millisecond)
			close(taskChannel)
			multiplier.Stop()

			assert.Equal(test, uint64(3), actionCount.Load())
			goleak.VerifyNone(test)
		},
	)

	test.Run(
		"Create multiplier with 10 workers count. Stop multiplier twice. On second call of stop execution control return immediately cause pool is already stopped.",
		func(test *testing.T) {
			thread := func() {
				time.Sleep(100 * time.Millisecond)
			}

			multiplier := multiplier.New(thread).
				WithMultiplyInterval(10 * time.Millisecond).
				Run()

			multiplier.Stop()
			multiplier.Stop()

			goleak.VerifyNone(test)
		},
	)
}
