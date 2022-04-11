package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type Result struct {
	Url          string
	Destinations []string
	Error        error
}

func (r Result) MarshalJSON() ([]byte, error) {
	type rs struct {
		Url          string   `json:"url"`
		Destinations []string `json:"destinations"`
		Error        string   `json:"error"`
	}
	n := rs{
		Url:          r.Url,
		Destinations: r.Destinations,
	}
	if r.Error != nil {
		n.Error = r.Error.Error()
	}
	var s bytes.Buffer
	enc := json.NewEncoder(&s)
	enc.SetEscapeHTML(false)
	err := enc.Encode(n)
	return s.Bytes(), err
}

type formatter interface {
	Start(url string) string
	Result(url string, destinations []string, err error) string
	Complete(results []Result) string
}

func errorToString(url string, err error) string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf(`Error dumping domain on "%s": `, url))
	if err == context.Canceled {
		s.WriteString("Operation aborted by the user")
	} else {
		s.WriteString(fmt.Sprintf(`%v`, err))
	}
	return s.String()
}
