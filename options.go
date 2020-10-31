package henge

type ConverterOpts struct {
	stringOpts
	mapOpts
}

type stringOpts struct {
	fmt  byte
	prec int
}

type mapOpts struct {
	maxDepth uint
}

func defaultConverterOpts() ConverterOpts {
	return ConverterOpts{
		stringOpts: stringOpts{
			fmt:  'f',
			prec: -1,
		},
		mapOpts: mapOpts{
			maxDepth: ^uint(0),
		},
	}
}

func withOpts(opts ConverterOpts) func(*ConverterOpts) {
	return func(dstOpts *ConverterOpts) {
		*dstOpts = opts
	}
}

func WithFloatFormat(fmt byte, prec int) func(*ConverterOpts) {
	return func(opt *ConverterOpts) {
		opt.stringOpts.fmt = fmt
		opt.stringOpts.prec = prec
	}
}

func WithMaxDepth(maxDepth uint) func(*ConverterOpts) {
	if maxDepth == 0 {
		panic("WithMaxDepth does not support zero")
	}
	return func(opt *ConverterOpts) {
		opt.mapOpts.maxDepth = maxDepth
	}
}
