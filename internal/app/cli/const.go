package cli

import "fmt"

type DialogMode int

const (
	WithoutContext DialogMode = 1
	WithContext    DialogMode = 2
)

func (d DialogMode) Name() string {
	switch d {
	case WithoutContext:
		return "without context"
	case WithContext:
		return "with context"
	default:
		return fmt.Sprintf("DialogMode(%d)", d)
	}
}
