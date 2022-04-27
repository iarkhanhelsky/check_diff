package executors

import (
	"go.uber.org/multierr"
	"sync"
)

var _ Executor = &Parallel{}

type Parallel struct {
	jobs []Job
}

func (exe *Parallel) Add(job Job) {
	exe.jobs = append(exe.jobs, job)
}

func (exe *Parallel) Run() error {
	var wg sync.WaitGroup
	wg.Add(len(exe.jobs))
	errchan := make(chan error)
	done := make(chan bool)
	for i, _ := range exe.jobs {
		job := exe.jobs[i]
		go func() {
			errchan <- job()
			wg.Done()
		}()
	}

	var err error
	go func() {
		for {
			select {
			case <-done:
				return
			case e := <-errchan:
				err = multierr.Append(err, e)
			}
		}
	}()

	wg.Wait()
	done <- true

	return err
}
