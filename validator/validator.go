package validator

// Check https://github.com/asaskevich/govalidator/blob/master/validator.go
// for more validation method if you need to implement others

type Validator interface {
	Validate() error
}

type Fn func() error

func Validate(fns ...Fn) (err error) {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return
}
