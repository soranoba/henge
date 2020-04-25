package henge

import "reflect"

// An element returned by Iterators
type Pair struct {
	Key   reflect.Value
	Value reflect.Value
}

// Iterator is a helper that sequentially accessing Array, Slice, Map, or Struct elements.
type Iterator interface {
	// Returns the next element, or nil if the end of elements has been reached.
	More() *Pair
	// Returns the number of accessible elements
	Count() int
}

type arrayIterator struct {
	value reflect.Value
	idx   int
	count int
}

type structIterator struct {
	value  reflect.Value
	fields []reflect.StructField
	idx    int
}

type mapIterator struct {
	value reflect.Value
	idx   int
	keys  []reflect.Value
}

// NewMustIterator is the same as NewIterator.
// If `reflect.Value` is not supported, it panics.
func NewMustIterator(v reflect.Value) Iterator {
	ite, ok := NewIterator(v)
	if !ok {
		panic("failed to create an iterator")
	}
	return ite
}

// Create an Iterator that sequentially accesses Array, Slice, Map, or Struct elements.
// If `reflect.Value` is not supported, it returns (nil, false).
func NewIterator(v reflect.Value) (Iterator, bool) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Array:
		return &arrayIterator{
			value: v,
			count: v.Len(),
		}, true
	case reflect.Slice:
		return &arrayIterator{
			value: v,
			count: v.Len(),
		}, true
	case reflect.Map:
		return &mapIterator{
			value: v,
			keys:  v.MapKeys(),
		}, true
	case reflect.Struct:
		return &structIterator{
			value:  v,
			fields: getFields(v),
		}, true
	default:
		return nil, false
	}
}

// Returns the next element, or nil if the end of elements has been reached.
func (ite *arrayIterator) More() *Pair {
	if ite.idx < ite.count {
		k := ite.idx
		v := ite.value.Index(ite.idx)
		ite.idx += 1
		return &Pair{Key: reflect.ValueOf(k), Value: v}
	}
	return nil
}

// Returns the number of accessible elements
func (ite *arrayIterator) Count() int {
	return ite.count
}

// Returns the next element, or nil if the end of elements has been reached.
func (ite *mapIterator) More() *Pair {
	if ite.idx < len(ite.keys) {
		k := ite.keys[ite.idx]
		v := ite.value.MapIndex(k)
		ite.idx += 1
		return &Pair{Key: k, Value: v}
	}
	return nil
}

// Returns the number of accessible elements
func (ite *mapIterator) Count() int {
	return len(ite.keys)
}

// Returns the next element, or nil if the end of elements has been reached.
func (ite *structIterator) More() *Pair {
	if ite.idx < len(ite.fields) {
		field := ite.fields[ite.idx]
		v := ite.value.Field(ite.idx)
		ite.idx += 1
		return &Pair{Key: reflect.ValueOf(field.Name), Value: v}
	}
	return nil
}

// Returns the number of accessible elements
func (ite *structIterator) Count() int {
	return len(ite.fields)
}

func getFields(v reflect.Value) []reflect.StructField {
	t := v.Type()
	fields := make([]reflect.StructField, 0)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fields = append(fields, field)
	}
	return fields
}
