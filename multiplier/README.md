# Multiplier

**Описание:**\
Библиотека для запуска и поддержания worker в заданном количестве потоков.

**Парадигма:** многопоточность, ООП. 

**Область применения:**\
Продолжительные/циклические процессы с разнородным pipeline-ом, состоящим из этапов разной сложности.
Например, многопоточная постобработка событий из queue/databus 

**Мотивация:**\
Необходимость поддерживать непрерывную обработку данных в нескольких потоках.

**Принцип работы:**\
Принимает на вход worker/функцию, которая будет запущена в указанном количестве потоков.
При завершении функции внутри потока она снова запускается через указанный интервал времени.
Таким образом заданное количество потоков выполнения функции всегда будет стремиться к заданному значению.

**Особенности:**
- Panic должен быть обработана внутри функции. Не обработанный panic будет заглушен и функция запустится заново.
- Deadlock блокирует поток выполнения. Функция должна сама реализовать его обработку (timeout, context canceling, context.timerCtx, etc...).
  Не обработанный deadlock уменьшает количество потоков.

**Примеры использования:**\
Сервис принимает slice c данными, каждый элемент которого может быть обработан параллельно
```go
package service

import (
	"context"
	"time"
	
	"github.com/pkg/errors"

	"go.avito.ru/gl/logger/v3"
    "go.avito.ru/gl/multiplier"
)

type Service struct {
    logger logger.Logger
}

func(s *Service) DoStuff(ctx context.Context, data []string) {
    tasks := make(chan string)
    workers := multiplier.New(s.handler(ctx, tasks)).
        WithWorkerCount(5).
        WithMultiplyInterval(10 * time.Millisecond).
        Run()

    defer func() {
        close(tasks)
        workers.Stop()
    }()

    for i := range data {
        tasks <- data[i]
    }
}

func(s *Service) handler(ctx context.Context, tasks <-chan string) func() {
    return func() {
        timer := time.NewTimer(10*time.Second)
        defer func() {
			timer.Stop()
            if err := recover(); err != nil {
                s.logger.Error(
                    ctx,
					errors.Errorf("Got panic on handle due to %s", err),
                )
            }
        }()

        for {
            select {
            case task, isOpen := <-tasks:
                if !isOpen{
                    return
                }

                //do some stuff here
            case <-timer.C:
                return
            case <-ctx.Done():  
				return
            }
        }
    }
}
```

Сервис принимает slice c данными, каждый элемент которого может быть обработан параллельно.
Однако обработка состоит из двух независимых операций разной сложности/скорости. Поэтому для оптимизации она разделена на два этапа.
```go
package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"go.avito.ru/gl/logger/v3"
	"go.avito.ru/gl/multiplier"
)

type Service struct {
	logger logger.Logger
}

func (s *Service) DoStuff(ctx context.Context, data []string) {
	firstStep := make(chan string, len(data))
	secondStep := make(chan string, len(data))
	firstStepWorkers := multiplier.New(s.heavyHandler(ctx, firstStep, secondStep)).
        WithWorkerCount(50).
		Run()

	secondStepWorkers := multiplier.New(s.lightHandler(ctx, secondStep)).
        WithWorkerCount(5).
		Run()

	defer func() {
		close(firstStep)
		firstStepWorkers.Stop()
		close(secondStep)
		secondStepWorkers.Stop()
	}()

	for i := range data {
		firstStep <- data[i]
	}
}

func (s *Service) heavyHandler(ctx context.Context, firstStep <-chan string, secondStep chan<- string) func() {
	return func() {
		timer := time.NewTimer(5 * time.Second)
		defer func() {
            timer.Stop()
			if err := recover(); err != nil {
				s.logger.Error(
                    ctx,
					errors.Errorf("Got panic cause work is too much heavy or impassable due to %s", err),
                )
			}
		}()

		for {
			select {
			case task, isOpen := <-firstStep:
				if !isOpen {
					return
				}

				//do some really heavy things here
				secondStep <- "result of magic and lucky, that ready for next step"
			case <-timer.C:
				return
            case <-ctx.Done():
              return
			}
		}
	}
}

func (s *Service) lightHandler(ctx context.Context, secondStep <-chan string) func() {
	return func() {
		timer := time.NewTimer(5 * time.Second)
		defer func() {
            timer.Stop()
			if err := recover(); err != nil {
				s.logger.Error(
                    ctx,
					errors.Errorf("Got panic cause paws are sweaty due to %s", err),
				)
			}
		}()

		for {
			select {
			case task, isOpen := <-secondStep:
				if !isOpen {
					return
				}

				//do some funny here, cause we can
			case <-timer.C:
				return
            case <-ctx.Done():
				return
			}
		}
	}
}
```