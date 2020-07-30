package errwrp

import "strings"

type ErrorPkg int

const (
	PkgErrors ErrorPkg = iota
	XErrors
	Fmt
	Unknown
)

func (e ErrorPkg) PkgName() string {
	switch e {
	case PkgErrors:
		return "github.com/pkg/errors"
	case XErrors:
		return "golang.org/x/xerrors"
	case Fmt:
		return "fmt"
	case Unknown:
		return "unknown"
	}
	panic("unreachable")
}

func FromPkgName(name string) ErrorPkg {
	n := strings.Trim(name, "\"")
	if n == PkgErrors.PkgName() {
		return PkgErrors
	}
	if n == XErrors.PkgName() {
		return XErrors
	}
	if n == Fmt.PkgName() {
		return Fmt
	}
	return Unknown
}
