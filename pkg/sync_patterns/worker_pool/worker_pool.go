package workerpool

import (
	"context"
	"sync"
)

////////////////////// Worker //////////////////////
type Worker[InType, OutType any] struct {
	inputCh  <-chan InType
	resultCh chan<- WorkerResult[OutType]
}

func NewWorker[InType, OutType any](inputCh <-chan InType, resultCh chan<- WorkerResult[OutType]) *Worker[InType, OutType] {
	return &Worker[InType, OutType]{
		inputCh:  inputCh,
		resultCh: resultCh,
	}
}

// запускает воркер
func (w *Worker[InType, OutType]) Run(ctx context.Context, wg *sync.WaitGroup, payload func(ctx context.Context, item InType) (OutType, error)) {
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case param, ok := <-w.inputCh:
				if !ok {
					return
				}
				res, err := payload(ctx, param) // выполнение основного действия
				if err != nil {
					w.resultCh <- WorkerResult[OutType]{Err: err}
					return
				}
				w.resultCh <- WorkerResult[OutType]{Result: res}
			}
		}
	}()
}

// структура для хранения результата выполнения задания воркером
type WorkerResult[OutType any] struct {
	Result OutType
	Err    error
}

////////////////////// WorkerPool //////////////////////
type WorkerPool[InType, OutType any] struct {
	workersCount  int
	workers       []*Worker[InType, OutType]
	input         <-chan InType              // входной канал на чтение или созданный канал на основе входного слайса
	result        chan WorkerResult[OutType]
	ctxCancelFunc context.CancelFunc
}

func NewWorkerPoolForData[InType, OutType any](workersCount int, data []InType) *WorkerPool[InType, OutType] {
	items := make(chan InType)

	go func() {
		for _, item := range data {
			items <- item
		}
		close(items)
	}()

	pool := &WorkerPool[InType, OutType]{
		workersCount: workersCount,
		input:  items,
		result: make(chan WorkerResult[OutType]),
	}

	pool.createWorkers()
	return pool
}

func NewWorkerPoolForChan[InType, OutType any](workersCount int, inputCh <-chan InType) *WorkerPool[InType, OutType] {
	pool := &WorkerPool[InType, OutType]{
		workersCount: workersCount,
		input:  inputCh,
		result: make(chan WorkerResult[OutType]),
	}

	pool.createWorkers()
	return pool
}

// создает воркеров и добавляет их в пул
func (p *WorkerPool[InType, OutType]) createWorkers() {
	for i := 0; i < p.workersCount; i++ {
		p.workers = append(p.workers, NewWorker[InType, OutType](p.input, p.result))
	}
}

// запускает воркеры
func (p *WorkerPool[InType, OutType]) Run(ctx context.Context, f func(context.Context, InType) (OutType, error)) <-chan WorkerResult[OutType] {
	ctx, cancel := context.WithCancel(ctx)
	p.ctxCancelFunc = cancel

	wg := sync.WaitGroup{}
	wg.Add(p.workersCount)

	for _, worker := range p.workers {
		worker.Run(ctx, &wg, f)
	}

	go func(wg *sync.WaitGroup){
		wg.Wait()
		close(p.result)
	}(&wg)

	return p.result
}

// останавливает воркеров
func (p *WorkerPool[InType, OutType]) Stop() {
	p.ctxCancelFunc()
}
