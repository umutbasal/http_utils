package http_range

import (
	"net/http"
	"sort"
)

// By is the type of a "less" function that defines the ordering of its http.Response arguments.
type By func(p1, p2 *http.Response) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(responses []*http.Response) {
	ps := &responseSorter{
		responses: responses,
		by:        by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// responseSorter joins a By function and a slice of http.Responses to be sorted.
type responseSorter struct {
	responses []*http.Response
	by        func(p1, p2 *http.Response) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *responseSorter) Len() int {
	return len(s.responses)
}

// Swap is part of sort.Interface.
func (s *responseSorter) Swap(i, j int) {
	s.responses[i], s.responses[j] = s.responses[j], s.responses[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *responseSorter) Less(i, j int) bool {
	return s.by(s.responses[i], s.responses[j])
}
