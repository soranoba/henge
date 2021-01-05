package henge

import "reflect"

type (
	// ConverterOpts are options for the conversion.
	ConverterOpts struct {
		stringOpts
		mapOpts
	}
	stringOpts struct {
		fmt  byte
		prec int
	}
	mapOpts struct {
		maxDepth   uint
		filterFuns mapFilterFuns
	}
	mapFilterFuns []func(k interface{}, v interface{}) bool
)

func (fs mapFilterFuns) All(k interface{}, v interface{}) bool {
	for _, f := range fs {
		if !f(k, v) {
			return false
		}
	}
	return true
}

func defaultConverterOpts() ConverterOpts {
	return ConverterOpts{
		stringOpts: stringOpts{
			fmt:  'f',
			prec: -1,
		},
		mapOpts: mapOpts{
			maxDepth:   ^uint(0),
			filterFuns: make(mapFilterFuns, 0),
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

// WithMapMaxDepth is an option when converting to map.
//
// By default, all structs are converted to maps.
// It can be used when converting only the top-level.
func WithMapMaxDepth(maxDepth uint) func(*ConverterOpts) {
	return func(opt *ConverterOpts) {
		opt.mapOpts.maxDepth = maxDepth
	}
}

// WithMapFilter is an option when converting to map.
//
// If you specify multiple filters, it will be copied only if all filters return true.
// By default, it copies everything.
func WithMapFilter(cond func(k interface{}, v interface{}) bool) func(*ConverterOpts) {
	return func(opt *ConverterOpts) {
		opt.mapOpts.filterFuns = append(opt.mapOpts.filterFuns, cond)
	}
}

// WithoutNilMapKey is an option when converting to map.
//
// When it used, it will not copy if the key is nil.
func WithoutNilMapKey() func(*ConverterOpts) {
	return WithMapFilter(func(k interface{}, v interface{}) bool {
		return !isNil(k)
	})
}

// WithoutNilMapValue is an option when converting to map.
//
// When it used, it will not copy if the value is nil.
func WithoutNilMapValue() func(*ConverterOpts) {
	return WithMapFilter(func(k interface{}, v interface{}) bool {
		return !isNil(v)
	})
}

// WithoutZeroMapKey is an option when converting to map.
//
// When it used, it will not copy if the key is zero.
func WithoutZeroMapKey() func(*ConverterOpts) {
	return WithMapFilter(func(k interface{}, v interface{}) bool {
		return !isZero(k)
	})
}

// WithoutZeroMapValue is an option when converting to map.
//
// When it used, it will not copy if the value is zero.
func WithoutZeroMapValue() func(*ConverterOpts) {
	return WithMapFilter(func(k interface{}, v interface{}) bool {
		return !isZero(v)
	})
}

func isNil(i interface{}) bool {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Interface, reflect.Slice, reflect.Chan, reflect.Func,
		reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		return v.IsNil()
	default:
		return false
	}
}

func isZero(i interface{}) bool {
	v := reflect.ValueOf(i)
	return !v.IsValid() || v.IsZero()
}
