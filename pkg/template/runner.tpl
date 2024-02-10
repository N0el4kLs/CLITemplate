package runner

type Runner struct {
    options *Options
}

func NewRunner(option *Options) (*Runner, error) {
	runner := &Runner{
	    options: option,
	}

	return runner, nil
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Run() {

}

func (r *Runner) Close() {

}