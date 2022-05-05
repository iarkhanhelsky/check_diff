package executors

type Executor interface {
	Add(job Job)
	Run() error
}

type Job func() error

func NewParallel() Executor {
	return &Parallel{}
}

func NewSequential() Executor {
	return &Sequential{}
}
