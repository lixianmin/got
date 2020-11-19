package loom

type options struct {
	Size      int
	CloseChan chan struct{}
}

type Option func(*options)

func createOptions(optionList []Option) options {
	var opts = options{
		Size: 8,
	}

	for _, opt := range optionList {
		opt(&opts)
	}

	if opts.CloseChan == nil {
		opts.CloseChan = make(chan struct{})
	}

	return opts
}

func WithSize(size int) Option {
	return func(opts *options) {
		if size > 0 {
			opts.Size = size
		}
	}
}

func WithCloseChan(c chan struct{}) Option {
	return func(opts *options) {
		if c != nil {
			opts.CloseChan = c
		}
	}
}
