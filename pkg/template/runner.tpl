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

func (r *Runner) Run() error {

    return nil
}

func (r *Runner) Close() {

}