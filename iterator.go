package henge

import "reflect"

type Pair struct {
	Key   reflect.Value
	Value reflect.Value
}

type Iterator interface {
	More() *Pair
	Count() int
}

type arrayIterator struct {
	value reflect.Value
	idx   int
	count int
}

type structIterator struct {
	value reflect.Value
	idx   int
	count int
}

type mapIterator struct {
	value reflect.Value
	idx   int
	keys  []reflect.Value
}

type emptyIterator struct {
}

func NewMustIterator(v reflect.Value) Iterator {
	ite, ok := NewIterator(v)
	if !ok {
		panic("failed to create an iterator")
	}
	return ite
}

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
			value: v,
			count: v.NumField(),
		}, true
	default:
		return nil, false
	}
}

func (ite *arrayIterator) More() *Pair {
	if ite.idx < ite.count {
		k := ite.idx
		v := ite.value.Index(ite.idx)
		ite.idx += 1
		return &Pair{Key: reflect.ValueOf(k), Value: v}
	}
	return nil
}

func (ite *arrayIterator) Count() int {
	return ite.count
}

func (ite *mapIterator) More() *Pair {
	if ite.idx < len(ite.keys) {
		k := ite.keys[ite.idx]
		v := ite.value.MapIndex(k)
		ite.idx += 1
		return &Pair{Key: k, Value: v}
	}
	return nil
}

func (ite *mapIterator) Count() int {
	return len(ite.keys)
}

func (ite *structIterator) More() *Pair {
	if ite.idx < ite.count {
		field := ite.value.Type().Field(ite.idx)
		v := ite.value.Field(ite.idx)
		ite.idx += 1
		return &Pair{Key: reflect.ValueOf(field.Name), Value: v}
	}
	return nil
}

func (ite *structIterator) Count() int {
	return ite.count
}
