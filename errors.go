package lenspath

import "fmt"

type EmptyLensPathErr struct{}

func (e *EmptyLensPathErr) Error() string {
	return "lenspath: must have at least one lens"
}

type InvalidLensPathErr struct {
	index int
}

func (e *InvalidLensPathErr) Error() string {
	return fmt.Sprintf("lenspath: could not navigate further, end of structure reached at index %d of lens", e.index)
}

func (e *InvalidLensPathErr) Is(err error) bool {
	_, ok := err.(*InvalidLensPathErr)
	return ok
}
