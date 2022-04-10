package splitdumper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

//Dump dumps a split.to link with given list of Option
//Url is the split.to link to dump, limit is the minimum requirement for each links to be seen
//returns a list of destination links
//presence of error denotes if links should be considered partial/incomplete
//When given context gets cancelled, error will be context.Canceled
//If no client has been provided, default http client will be used with a timeout
//this behaviour should not be relied upon, using WithClient is recommended.
func Dump(url string, limit int, opts ...Option) ([]string, error) {
	c := config{
		url:     url,
		limit:   limit,
		context: context.Background(),
	}
	for _, o := range opts {
		o(&c)
	}

	if c.client == nil {
		hc := *http.DefaultClient
		const timeout = 10 * time.Second
		hc.Timeout = timeout
		c.client = &hc
	}
	if c.wait == nil {
		c.wait = func() {}
	}
	if c.hook == nil {
		c.hook = func(_ string, _ int) {}
	}
	return dump(c)
}

//dump is the actual implementation of Dump
func dump(c config) ([]string, error) {
	var err error

	r, err2 := http.NewRequestWithContext(c.context, "HEAD", c.url, nil)
	if err2 != nil {
		return nil, fmt.Errorf("error creating request: %w", err2)
	}
	seen := make(map[string]int)
loop:
	for {
		res, err3 := c.client.Do(r)
		if err3 != nil {
			err = err3
			switch {
			case errors.Is(err3, context.Canceled):
				err = context.Canceled
			}
			break
		}
		dest := res.Request.URL.String()
		_ = res.Body.Close()
		if dest == "" {
			err = errors.New(`location is empty`)
			break
		}

		seen[dest]++

		c.hook(dest, seen[dest])

		stop := true
		for _, count := range seen {
			if count < c.limit {
				stop = false
				break
			}
		}
		if stop {
			break
		}
		select {
		case <-c.context.Done():
			err = c.context.Err()
			break loop
		default:
		}
		c.wait()
	}
	urls := make([]string, 0, len(seen))
	for url := range seen {
		urls = append(urls, url)
	}

	return urls, err
}
