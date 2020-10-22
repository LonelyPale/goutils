package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

func UnknownError(r interface{}) error {
	switch e := r.(type) {
	case error:
		return e
	case nil:
		return nil
	case string:
		return errors.New(e)
	default:
		return errors.New(fmt.Sprint(r))
	}
}

func Errors(errs ...error) error {
	es := make([]error, 0)
	for _, e := range errs {
		if e != nil {
			es = append(es, e)
		}
	}

	switch len(es) {
	case 0:
		return nil
	case 1:
		return es[0]
	default:
		err := errors.New("errors:")
		for _, e := range es {
			err = Wrap(err, e.Error())
		}
		return err
	}
}

// import github.com/pkg/errors

func New(message string) error {
	return errors.New(message)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func WithStack(err error) error {
	return errors.WithStack(err)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

func WithMessagef(err error, format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args...)
}

func Cause(err error) error {
	return errors.Cause(err)
}
