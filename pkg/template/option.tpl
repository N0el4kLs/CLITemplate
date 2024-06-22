package runner

type Options struct {
}

func ParseOptions() (*Options,error) {
    ShowBanner()

	options := &Options{}
	return options,nil
}