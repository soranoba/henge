package henge

import (
	"reflect"
)

const (
	structTagKey = "henge"
)

type structTag struct {
	ignore bool
}

// newStructTag is create a `structTag` from `reflect.StructField`
func newStructTag(f reflect.StructField) structTag {
	value := f.Tag.Get(structTagKey)

	switch value {
	case "-":
		return structTag{ignore: true}
	default:
		return structTag{ignore: false}
	}
}

// getStructFieldIndexes returns all field indexes including embedded fields of the type.
func getStructFieldIndexes(t reflect.Type) [][]int {
	fieldIndexes := make([][]int, 0)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fieldIndexes = append(fieldIndexes, f.Index)
			if f.Anonymous {
				for _, index := range getStructFieldIndexes(f.Type) {
					fieldIndexes = append(fieldIndexes, append(f.Index, index...))
				}
			}
		}
	}
	return fieldIndexes
}

type structField struct {
	name  string
	index []int
	tags  []structTag
}

// getStructFields returns all fields including embedded fields of the type.
func getStructFields(t reflect.Type) []structField {
	fieldIndexes := getStructFieldIndexes(t)
	fields := make([]structField, len(fieldIndexes))

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		for i, fieldIndex := range fieldIndexes {
			f := t.FieldByIndex(fieldIndex)
			tags := make([]structTag, len(fieldIndex))
			for i := 0; i < len(fieldIndex); i++ {
				tags[i] = newStructTag(t.FieldByIndex(fieldIndex[0 : i+1]))
			}
			fields[i] = structField{
				name:  f.Name,
				index: fieldIndex,
				tags:  tags,
			}
		}
	}
	return fields
}

func (f *structField) isIgnore() bool {
	for _, t := range f.tags {
		if t.ignore {
			return true
		}
	}
	return false
}

func (f *structField) isMatch(rf reflect.StructField) bool {
	if len(f.index) != len(rf.Index) {
		return false
	}
	if f.name != rf.Name {
		return false
	}
	for i := 0; i < len(f.index); i++ {
		if f.index[i] != rf.Index[i] {
			return false
		}
	}
	return true
}
