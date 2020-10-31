package henge

type converter struct {
	isNil   bool
	field   string
	opts    ConverterOpts
	storage map[string]interface{}
}

func (c *converter) new(i interface{}, fieldName string) *ValueConverter {
	newConverter := New(i)
	newConverter.converter = *c
	newConverter.converter.field = fieldName
	return newConverter
}

func (c *converter) InstanceGet(key string) (interface{}, bool) {
	v, ok := c.storage[key]
	return v, ok
}

func (c *converter) InstanceSet(key string, value interface{}) {
	c.storage[key] = value
}
