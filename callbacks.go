package henge

// BeforeCallback is a callback that is executed before the conversion from a struct to the struct.
// Deprecated: Use BeforeConvertFromCallback instead.
type BeforeCallback interface {
	BeforeConvert(src interface{}, store InstanceStore) error
}

// AfterCallback is a callback that is executed after the conversion from a struct to the struct.
// Deprecated: Use AfterConvertFromCallback instead.
type AfterCallback interface {
	AfterConvert(src interface{}, store InstanceStore) error
}

// BeforeConvertFromCallback is a callback that is executed before the conversion from a struct to the struct.
// To define structures across packages, you need to define them in the package's struct that imports other packages.
// Within structures of the same package, you must use BeforeConvertFromCallback. It cannot be used simultaneously with BeforeConvertToCallback.
type BeforeConvertFromCallback interface {
	BeforeConvertFrom(src interface{}, store InstanceStore) error
}

// AfterConvertFromCallback is a callback that is executed after the conversion from a struct to the struct.
// To define structures across packages, you need to define them in the package's struct that imports other packages.
// Within structures of the same package, you must use AfterConvertFromCallback. It cannot be used simultaneously with AfterConvertToCallback.
type AfterConvertFromCallback interface {
	AfterConvertFrom(src interface{}, store InstanceStore) error
}

// BeforeConvertToCallback is a callback that is executed before the conversion from the struct to a struct.
// To define structures across packages, you need to define them in the package's struct that imports other packages.
// Within structures of the same package, you must use BeforeConvertFromCallback. It cannot be used simultaneously with BeforeConvertToCallback.
type BeforeConvertToCallback interface {
	BeforeConvertTo(dst interface{}, store InstanceStore) error
}

// AfterConvertToCallback is a callback that is executed after the conversion from the struct to a struct.
// To define structures across packages, you need to define them in the package's struct that imports other packages.
// Within structures of the same package, you must use AfterConvertFromCallback. It cannot be used simultaneously with AfterConvertToCallback.
type AfterConvertToCallback interface {
	AfterConvertTo(dst interface{}, store InstanceStore) error
}
