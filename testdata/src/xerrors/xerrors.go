package xerrors

import "golang.org/x/xerrors"

func method() error {
	return xerrors.New("example")
}

func method2() error {
	return method() // want `should use errors\.Wrap\(\) or errors\.Wrapf\(\)`
}
