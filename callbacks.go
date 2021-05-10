package henge

// BeforeCallback is a callback that is executed before the conversion from a struct to the struct.
type BeforeCallback interface {
	BeforeConvert(src interface{}, store InstanceStore) error
}

// AfterCallback is a callback that is executed after the conversion from a struct to the struct.
type AfterCallback interface {
	AfterConvert(src interface{}, store InstanceStore) error
}
