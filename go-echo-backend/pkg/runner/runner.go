package runner

import (
	"sync"

	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/panjf2000/ants/v2"
)

type Runner struct {
	wg     *sync.WaitGroup
	pool   *ants.Pool
	logger *Logger
}

type DoFunc = func(i interface{})

// New creates a new job runner.
func New(workers int, preAlloc ...bool) *Runner {
	var wg = &sync.WaitGroup{}
	var logger = &Logger{
		Logger: logger.New("runner"),
	}
	var panicHandler = func(i interface{}) {
		logger.ErrorAny(i, "Panic occured")
	}

	var validWorker = 1
	if workers > 0 {
		validWorker = workers
	}

	pool, err := ants.NewPool(validWorker, func(opts *ants.Options) {
		opts.Logger = logger
		opts.PanicHandler = panicHandler

		if len(preAlloc) > 0 {
			opts.PreAlloc = preAlloc[0]
		}
	})
	if err != nil {
		panic(err)
	}

	return &Runner{
		wg:     wg,
		pool:   pool,
		logger: logger,
	}
}

func (r *Runner) Submit(task func()) error {
	r.wg.Add(1)
	return r.pool.Submit(func() {
		defer r.wg.Done()
		task()
	})
}

func (r *Runner) Wait() {
	r.wg.Wait()
}

func (r *Runner) Release() {
	r.pool.Release()
}
