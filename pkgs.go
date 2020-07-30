package errwrp

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
	if name == PkgErrors.PkgName() {
		return PkgErrors
	}
	if name == XErrors.PkgName() {
		return XErrors
	}
	if name == Fmt.PkgName() {
		return Fmt
	}
	return Unknown
}