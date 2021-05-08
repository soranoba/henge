package henge

import (
	"math"
	"reflect"
)

type (
	// ConverterOption are options for the conversion.
	ConverterOption func(opts *converterOpts)
	// RoundingFunc is a function that rounds from float to nearest integer.
	// e.g. math.Floor
	RoundingFunc func(float64) float64
	// ConversionFunc is a conversion function of values.
	ConversionFunc func(converter *ValueConverter) Converter
)

var (
	// DefaultConversionFunc is a ConversionFunc used by default.
	// It is no conversion.
	DefaultConversionFunc ConversionFunc = func(converter *ValueConverter) Converter {
		return converter
	}
)

type (
	converterOpts struct {
		numOpts
		stringOpts
		sliceOpts
		mapOpts
	}
	numOpts struct {
		roundingFunc RoundingFunc
	}
	stringOpts struct {
		fmt  byte
		prec int
	}
	sliceOpts struct {
		valueConversionFunc ConversionFunc
	}
	mapOpts struct {
		maxDepth            uint
		filterFuns          mapFilterFuns
		keyType             reflect.Type
		keyConversionFunc   ConversionFunc
		valueConversionFunc ConversionFunc
	}
	mapFilterFuns []func(k interface{}, v interface{}) bool
)

var (
	interfaceType = reflect.ValueOf([]interface{}{}).Type().Elem()
)

func (fs mapFilterFuns) All(k interface{}, v interface{}) bool {
	for _, f := range fs {
		if !f(k, v) {
			return false
		}
	}
	return true
}

func defaultConverterOpts() *converterOpts {
	return &converterOpts{
		numOpts: numOpts{
			roundingFunc: math.Floor,
		},
		stringOpts: stringOpts{
			fmt:  'f',
			prec: -1,
		},
		sliceOpts: sliceOpts{
			valueConversionFunc: DefaultConversionFunc,
		},
		mapOpts: mapOpts{
			maxDepth:            ^uint(0),
			filterFuns:          make(mapFilterFuns, 0),
			keyType:             interfaceType,
			keyConversionFunc:   DefaultConversionFunc,
			valueConversionFunc: DefaultConversionFunc,
		},
	}
}

// WithFloatFormat is an option when converting from float to string.
// Ref: strconv.FormatFloat
func WithFloatFormat(fmt byte, prec int) ConverterOption {
	return func(opt *converterOpts) {
		opt.stringOpts.fmt = fmt
		opt.stringOpts.prec = prec
	}
}

// WithRoundingFunc is an option when converting from float to integer (or unsigned integer).
// It specify the rounding method from float to nearest integer.
// By default, it use math.Floor.
func WithRoundingFunc(f RoundingFunc) ConverterOption {
	return func(opt *converterOpts) {
		opt.numOpts.roundingFunc = f
	}
}

// WithSliceValueConverter is an option when converting to slice.
//
// It can be used when converting values to other types.
func WithSliceValueConverter(f ConversionFunc) ConverterOption {
	return func(opt *converterOpts) {
		opt.sliceOpts.valueConversionFunc = f
	}
}

// WithMapKeyConverter is an option when converting to map.
//
// It can be used when converting keys to other types.
func WithMapKeyConverter(f ConversionFunc) ConverterOption {
	return func(opt *converterOpts) {
		opt.mapOpts.keyConversionFunc = f
	}
}

// WithMapValueConverter is an option when converting to map.
//
// It can be used when converting values to other types.
func WithMapValueConverter(f ConversionFunc) ConverterOption {
	return func(opt *converterOpts) {
		opt.mapOpts.valueConversionFunc = f
	}
}

// WithMapMaxDepth is an option when converting to map.
//
// By default, all structs are converted to maps.
// It can be used when converting only the top-level.
func WithMapMaxDepth(maxDepth uint) ConverterOption {
	return func(opt *converterOpts) {
		opt.mapOpts.maxDepth = maxDepth
	}
}

// WithMapFilter is an option when converting to map.
//
// If you specify multiple filters, it will be copied only if all filters return true.
// By default, it copies everything.
func WithMapFilter(cond func(k interface{}, v interface{}) bool) ConverterOption {
	return func(opt *converterOpts) {
		opt.mapOpts.filterFuns = append(opt.mapOpts.filterFuns, cond)
	}
}

// WithoutNilMapKey is an option when converting to map.
//
// When it used, it will not copy if the key is nil.
func WithoutNilMapKey() ConverterOption {
	return WithMapFilter(func(k interface{}, v interface{}) bool {
		return !isNil(k)
	})
}

// WithoutNilMapValue is an option when converting to map.
//
// When it used, it will not copy if the value is nil.
func WithoutNilMapValue() ConverterOption {
	return WithMapFilter(func(k interface{}, v interface{}) bool {
		return !isNil(v)
	})
}

// WithoutZeroMapKey is an option when converting to map.
//
// When it used, it will not copy if the key is zero.
func WithoutZeroMapKey() ConverterOption {
	return WithMapFilter(func(k interface{}, v interface{}) bool {
		return !isZero(k)
	})
}

// WithoutZeroMapValue is an option when converting to map.
//
// When it used, it will not copy if the value is zero.
func WithoutZeroMapValue() ConverterOption {
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
