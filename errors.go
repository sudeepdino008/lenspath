package lenspath

import "fmt"

type EmptyLensPathErr struct{}

func (e *EmptyLensPathErr) Error() string {
	return "lenspath: must have at least one lens"
}

type InvalidLensPathErr struct {
	index   int
	errType InvalidLensPathErrType
}

type InvalidLensPathErrType string

const (
	ArrayExpectedErr   InvalidLensPathErrType = "lenspath: expected array (*)"
	LensPathStoppedErr                        = "lenspath: could not navigate further, end of structure reached"
)

func NewInvalidLensPathErr(index int, errType InvalidLensPathErrType) *InvalidLensPathErr {
	return &InvalidLensPathErr{index, errType}
}

func (e *InvalidLensPathErr) Error() string {
	return fmt.Sprintf("lenspath: %s; lens index: %d", e.errType, e.index)
}

func (e *InvalidLensPathErr) Is(err error) bool {
	_, ok := err.(*InvalidLensPathErr)
	return ok
}
