package main

import "strings"

var _ formatter = (*rawFormat)(nil)

type rawFormat struct{}

func (r rawFormat) Start(string) string {
	return ""
}

func (r rawFormat) Result(url string, destinations []string, err error) string {
	var s strings.Builder
	for i, dest := range destinations {
		if i > 0 {
			s.WriteString("\n")
		}
		s.WriteString(dest)
	}

	if err != nil {
		s.WriteString("\n")
		s.WriteString(errorToString(url, err))
	}
	return s.String()
}

func (r rawFormat) Complete([]Result) string {
	return ""
}
