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
	ArrayExpectedErr         InvalidLensPathErrType = "expected array (*)"
	LensPathStoppedErr                              = "could not navigate further, end of structure reached"
	CannotSetFieldErr                               = "cannot set field"
	PathContainsArrErr                              = "path contains *; use GetList() instead"
	PathDoesntContainsArrErr                        = "path does not contain *; use Get() instead"
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

type InvalidSetParamErr string

const (
	ArrayParamExpectedErr InvalidSetParamErr = "expected array for set value"
	ParamSizeMismatchErr  InvalidSetParamErr = "array param and structure field array length should match"
)

func (e InvalidSetParamErr) Error() string {
	return fmt.Sprintf("lenspath: %s", string(e))
}
