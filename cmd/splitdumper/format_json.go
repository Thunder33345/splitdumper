package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

var _ formatter = (*jsonFormat)(nil)

type jsonFormat struct {
}

func (j jsonFormat) Start(string) string {
	return ""
}

func (j jsonFormat) Result(string, []string, error) string {
	return ""
}

func (j jsonFormat) Complete(results []Result) string {
	var s strings.Builder
	enc := json.NewEncoder(&s)
	enc.SetEscapeHTML(false)
	err := enc.Encode(results)
	if err != nil {
		return fmt.Sprintf("Failed to create JSON object %s", err)
	}
	return strings.TrimRight(s.String(), "\n")
}
