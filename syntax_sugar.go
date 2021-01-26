package henge

// ToString is equiv to New(i, fs...).String().Value()
func ToString(i interface{}, fs ...func(*ConverterOpts)) string {
	return New(i, fs...).String().Value()
}

// ToInt is equiv to New(i, fs...).Int().Value()
func ToInt(i interface{}, fs ...func(*ConverterOpts)) int64 {
	return New(i, fs...).Int().Value()
}

// ToUint is equiv to New(i, fs...).Uint().Value()
func ToUint(i interface{}, fs ...func(*ConverterOpts)) uint64 {
	return New(i, fs...).Uint().Value()
}

// ToFloat is equiv to New(i, fs...).Float().Value()
func ToFloat(i interface{}, fs ...func(*ConverterOpts)) float64 {
	return New(i, fs...).Float().Value()
}

// ToStringPtr is equiv to New(i, fs...).StringPtr().Value()
func ToStringPtr(i interface{}, fs ...func(*ConverterOpts)) *string {
	return New(i, fs...).StringPtr().Value()
}

// ToIntPtr is equiv to New(i, fs...).IntPtr().Value()
func ToIntPtr(i interface{}, fs ...func(*ConverterOpts)) *int64 {
	return New(i, fs...).IntPtr().Value()
}

// ToUintPtr is equiv to New(i, fs...).UintPtr().Value()
func ToUintPtr(i interface{}, fs ...func(*ConverterOpts)) *uint64 {
	return New(i, fs...).UintPtr().Value()
}

// ToFloatPtr is equiv to New(i, fs...).FloatPtr().Value()
func ToFloatPtr(i interface{}, fs ...func(*ConverterOpts)) *float64 {
	return New(i, fs...).FloatPtr().Value()
}
