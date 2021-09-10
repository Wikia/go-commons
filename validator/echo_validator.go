package validator

import (
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

type setupFunction func(validate *validator.Validate)

type EchoValidator struct {
	once     sync.Once
	validate *validator.Validate
	setupFn setupFunction
}

func (v *EchoValidator) Validate(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyInit()

		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

//SetSetupFunction allows validator customization. This function
// will be invoked when validator is being initialized for the very first time.
// It is best place to register custom validators
func (v *EchoValidator) SetSetupFunction(fn setupFunction) {
	v.setupFn = fn
}

func (v *EchoValidator) Engine() *EchoValidator {
	v.lazyInit()

	return v
}

func (v *EchoValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		// add any custom validations etc. here
		// see: https://pkg.go.dev/github.com/go-playground/validator/v10
		if v.setupFn != nil {
			v.setupFn(v.validate)
		}
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}

	return valueType
}
