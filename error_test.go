package henge

import "testing"

func TestConvertError(t *testing.T) {
	var _ error = &ConvertError{}
}
