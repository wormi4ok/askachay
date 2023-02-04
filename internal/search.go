package internal

import (
	"errors"
	"fmt"
	"strings"
)

var UnrecognizedInputErr = errors.New("unrecognized input")

type Searcher interface {
	SearchMusic(searchTerm string) (name, url string, err error)
}

type MultiSearcher struct {
	ss []Searcher
}

func NewMultiSearcher(ss ...Searcher) *MultiSearcher {
	return &MultiSearcher{ss: ss}
}

func (m *MultiSearcher) SearchMusic(searchTerm string) (name, url string, err error) {
	var searchErrs []string
	for _, s := range m.ss {
		name, url, err = s.SearchMusic(searchTerm)
		if err == nil {
			return name, url, err
		}
		if !errors.Is(err, UnrecognizedInputErr) {
			searchErrs = append(searchErrs, err.Error())
		}
	}

	if len(searchErrs) == 0 {
		return "", "", UnrecognizedInputErr

	}
	return "", "", fmt.Errorf(strings.Join(searchErrs, "\n"))
}
