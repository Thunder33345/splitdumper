package splitdumper

import (
	"errors"
	"net/http"
)

//Dump dumps a split.to link
//url is the split.to link to dump
//limit is how many times all known sites must be seen before stopping
func Dump(client http.Client, url string, limit int) ([]string, error) {
	return dump(client, url, limit, func() {})
}

//DumpWithWait dumps a split.to link while taking a blocking function that will be called at the end of every loop
func DumpWithWait(client http.Client, url string, limit int, wait func()) ([]string, error) {
	return dump(client, url, limit, wait)
}

func dump(client http.Client, url string, limit int, wait func()) ([]string, error) {
	seen := make(map[string]int)
	for {
		res, err := client.Head(url)
		if err != nil {
			return nil, err
		}
		dest := res.Request.URL.String()
		if dest == "" {
			return nil, errors.New(`location is empty`)
		}
		if _, ok := seen[dest]; ok {
			seen[dest]++
		} else {
			seen[dest] = 1
		}
		stop := true
		for _, count := range seen {
			if count < limit {
				stop = false
				break
			}
		}
		if stop {
			break
		}
		wait()
	}
	urls := make([]string, 0, len(seen))
	for url, _ := range seen {
		urls = append(urls, url)
	}
	return urls, nil
}
