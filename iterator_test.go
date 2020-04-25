package henge

import (
	"reflect"
	"testing"
)

func TestArrayIterator(t *testing.T) {
	in := []string{"a", "b", "c"}

	ite, ok := NewIterator(reflect.ValueOf(in))
	assertEqual(t, ok, true)
	assertEqual(t, ite.Count(), 3)

	got := allAccessablePairsMap(ite)
	assertEqual(t, got, map[interface{}]interface{}{
		0: "a",
		1: "b",
		2: "c",
	})
}

func TestSliceIterator(t *testing.T) {
	in := []string{"a", "b", "c"}[0:2]

	ite, ok := NewIterator(reflect.ValueOf(in))
	assertEqual(t, ok, true)
	assertEqual(t, ite.Count(), 2)

	got := allAccessablePairsMap(ite)
	assertEqual(t, got, map[interface{}]interface{}{
		0: "a",
		1: "b",
	})
}

func TestMapIterator(t *testing.T) {
	in := map[string]string{
		"a": "av",
		"b": "bv",
		"c": "cv",
	}

	ite, ok := NewIterator(reflect.ValueOf(in))
	assertEqual(t, ok, true)
	assertEqual(t, ite.Count(), 3)

	got := allAccessablePairsMap(ite)
	assertEqual(t, got, map[interface{}]interface{}{
		"a": "av",
		"b": "bv",
		"c": "cv",
	})
}

func TestStructIterator(t *testing.T) {
	in := struct {
		A string
		b string
	}{
		A: "a",
		b: "b",
	}

	ite, ok := NewIterator(reflect.ValueOf(in))
	assertEqual(t, ok, true)
	assertEqual(t, ite.Count(), 2)

	got := allAccessablePairsMap(ite)
	assertEqual(t, got, map[interface{}]interface{}{
		"A": "a",
	})
}

func allAccessablePairsMap(ite Iterator) map[interface{}]interface{} {
	got := make(map[interface{}]interface{})
	for {
		pair := ite.More()
		if pair == nil {
			break
		}
		// private field cannot be converted to interface
		if pair.Key.CanInterface() && pair.Value.CanInterface() {
			got[pair.Key.Interface()] = pair.Value.Interface()
		}
	}
	return got
}
