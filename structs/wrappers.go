package structs

type Comparable interface {
	Less(k Comparable) bool
}

// Type wrappers
type Int int
type Float float64

func (i Int) Less(j Comparable) bool {
	return i < j.(Int)
}

func (i Float) Less(j Comparable) bool {
	return i < j.(Float)
}
