package main

import (
	"fmt"
	"strings"
)

var _ formatter = (*textFormat)(nil)

type textFormat struct{}

func (t textFormat) Start(url string) string {
	return fmt.Sprintf(`Dumping urls from: %s`, url)
}

func (t textFormat) Result(url string, destinations []string, err error) string {
	var s strings.Builder
	if err != nil {
		s.WriteString(errorToString(url, err))
		s.WriteString("\n")
	}

	if err != nil {
		s.WriteString("(Incomplete)")
	}
	s.WriteString(fmt.Sprintf("Found %d destinations:", len(destinations)))
	for _, dest := range destinations {
		s.WriteString("\n")
		s.WriteString(dest)
	}

	return s.String()
}

func (t textFormat) Complete(results []Result) string {
	if len(results) <= 1 {
		return ""
	}
	total := 0
	for _, res := range results {
		total += len(res.Destinations)
	}
	return fmt.Sprintf("Dumped %d urls", total)
}
