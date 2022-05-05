package executors

var _ Executor = &Sequential{}

type Sequential struct {
	jobs []Job
}

func (exe *Sequential) Add(job Job) {
	exe.jobs = append(exe.jobs, job)
}

func (exe *Sequential) Run() error {
	for _, j := range exe.jobs {
		if err := j(); err != nil {
			return err
		}
	}

	return nil
}
