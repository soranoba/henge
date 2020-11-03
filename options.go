package henge

// ConverterOpts are options for the conversion.
type ConverterOpts struct {
	stringOpts
	mapOpts
}

type stringOpts struct {
	fmt  byte
	prec int
}

type mapOpts struct {
	maxDepth   uint
	filterFunc func(k interface{}, v interface{}) bool
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

// WithFloatFormat is an option when converting from float to string.
// Ref: strconv.FormatFloat
func WithFloatFormat(fmt byte, prec int) func(*ConverterOpts) {
	return func(opt *ConverterOpts) {
		opt.stringOpts.fmt = fmt
		opt.stringOpts.prec = prec
	}
}

// WithMaxDepth is an option when converting to map.
//
// By default, all structs are converted to maps.
// It can be used when converting only the top-level.
func WithMaxDepth(maxDepth uint) func(*ConverterOpts) {
	if maxDepth == 0 {
		panic("WithMaxDepth does not support zero")
	}
	return func(opt *ConverterOpts) {
		opt.mapOpts.maxDepth = maxDepth
	}
}

// WithFilter is an option when converting to map.
//
// By default, values is copied even if it is nil.
// You can use this option to prevent this.
func WithFilter(cond func(k interface{}, v interface{}) bool) func(*ConverterOpts) {
	return func(opt *ConverterOpts) {
		opt.mapOpts.filterFunc = cond
	}
}
