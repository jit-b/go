package multiplier

import (
	"sync/atomic"
	"time"
)

type multiplier struct {
	config              *config
	isRunning           atomic.Bool
	runningWorkersCount atomic.Uint64
}

func run(config *config) *multiplier {
	result := &multiplier{config: config}
	result.multiplyWorker()
	result.isRunning.Store(true)

	go result.updateWorkerCount()

	return result
}

func (m *multiplier) multiplyWorker() {
	for m.runningWorkersCount.Load() < m.config.workerCount.Load() {
		m.runningWorkersCount.Add(1)
		m.runNewWorker()
	}
}

func (m *multiplier) runNewWorker() {
	go func() {
		defer func() {
			_ = recover()
			m.runningWorkersCount.Add(^uint64(0))
		}()
		m.config.worker()
	}()
}

func (m *multiplier) updateWorkerCount() {
	ticker := time.NewTicker(m.config.workerMultiplyInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.multiplyWorker()
		default:
			if !m.isRunning.Load() {
				return
			}
		}
	}
}

// Stop блокирует вызов до момента завершения выполнения worker-a во всех потоках.
func (m *multiplier) Stop() {
	if !m.isRunning.Load() {
		return
	}
	m.isRunning.Store(false)

	m.config.workerCount.Store(0)
	for m.runningWorkersCount.Load() != 0 {
	}
}
