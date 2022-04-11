package splitdumper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	netUrl "net/url"
	"time"
)

//Dump dumps a split.to link with given list of Option
//Url is the split.to link to dump
//limit is the minimum times each destination needs to be seen
//use WithBreaker to define how the halting behavior
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
		breaker: ConservativeBreaker(),
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

	if !c.fullTrace {
		urlF, err := netUrl.Parse(c.url)
		if err != nil {
			return nil, fmt.Errorf("error parsing url: %w", err)
		}

		tmp := *c.client
		c.client = &tmp
		c.client.CheckRedirect = redirectChecker(urlF.Host)
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
	record := make(map[string]int)
loop:
	for {
		res, err3 := c.client.Do(r)
		if err3 != nil {
			err = err3
			if errors.Is(err3, c.context.Err()) {
				//handle error caused by context cancel or timeout
				err = c.context.Err()
			}
			break
		}
		_ = res.Body.Close()

		var dest string
		if !c.fullTrace {
			dest = res.Header.Get("Location")
		} else {
			dest = res.Request.URL.String()
		}

		if dest == "" {
			err = fmt.Errorf("returned destination is emtpy on %s", c.url)
			break
		}

		record[dest]++

		c.hook(dest, record[dest])

		if c.breaker(c.limit, dest, record) {
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
	urls := make([]string, 0, len(record))
	for url := range record {
		urls = append(urls, url)
	}

	return urls, err
}

func redirectChecker(host string) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if req.URL.Host == host && len(via) < 10 {
			return nil
		}
		return http.ErrUseLastResponse
	}
}
