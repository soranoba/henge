package henge

type BeforeCallback interface {
	BeforeConvert(src interface{}, converter Converter) error
}

type AfterCallback interface {
	AfterConvert(src interface{}, converter Converter) error
}
