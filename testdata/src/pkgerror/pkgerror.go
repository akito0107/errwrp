package pkgerror

import "github.com/pkg/errors"

func method() error {
	return errors.New("example")
}

func method2() error {
	return method() // want "should be use errors.Wrap() or errors.Wrapf()"
}
