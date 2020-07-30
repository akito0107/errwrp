package ignore

import "github.com/pkg/errors"

func method() error {
	return errors.New("example")
}

func method2() error {
	return method() //mustwrap: ignore
}

func method3() error {
	//mustwrap: ignore
	return method()
}
