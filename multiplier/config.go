package multiplier

import (
	"sync/atomic"
	"time"
)

const (
	defaultWorkersCount           = 1
	defaultMultiplyWorkerInterval = time.Second
)

type config struct {
	worker                 func()
	workerCount            atomic.Uint64
	workerMultiplyInterval time.Duration
}

// New создает конфигурацию запуска multiplier-a
func New(worker func()) *config {
	result := &config{
		worker:                 worker,
		workerMultiplyInterval: defaultMultiplyWorkerInterval,
	}

	return result.WithWorkerCount(defaultWorkersCount)
}

// WithWorkerCount устанавливает количество одновременно работающих worker-ов.
// По умолчанию worker выполняется в одном потоке.
func (c *config) WithWorkerCount(workerCount uint64) *config {
	if workerCount > 0 {
		c.workerCount.Store(workerCount)
	}

	return c
}

// WithMultiplyInterval устанавливает интервал запуска новых worker-ов (в случае если количество исполняемых меньше чем было установлено).
// По умолчанию проверка и запуск происходит каждую секунду.
func (c *config) WithMultiplyInterval(threadsCountUpdateInterval time.Duration) *config {
	if threadsCountUpdateInterval > 0 {
		c.workerMultiplyInterval = threadsCountUpdateInterval
	}

	return c
}

// Run запускает multiplier.
func (c *config) Run() *multiplier {
	return run(c)
}
