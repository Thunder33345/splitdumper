package main

import (
	"fmt"
	"github.com/Thunder33345/splitdumper"
	"github.com/jessevdk/go-flags"
	"net/http"
	"os"
	"sort"
	"time"
)

var opts struct {
	Timeout time.Duration `short:"t" long:"timeout" description:"Sets the client timeout" default:"5s"`
	Wait    time.Duration `short:"w" long:"wait" description:"Sets the wait time after crawling"  default:"50ms"`
	Limit   int           `short:"l" long:"limit" description:"How many times all url should be seen before stopping" default:"3"`
	Raw     bool          `short:"r" long:"raw" description:"Outputs only the end result"`
	Args    struct {
		Urls []string `description:"The urls to dump (required)" required:"1"`
	} `positional-args:"yes"`
}

func main() {
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	_, err := parser.Parse()
	if flags.WroteHelp(err) {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
		return
	}

	if err != nil {
		fmt.Printf("Error parsing arguments: %v\nUse %v --help to see help", err, os.Args[0])
		os.Exit(1)
		return
	}

	if len(opts.Args.Urls) <= 0 {
		fmt.Printf("Missing urls to dump")
		os.Exit(1)
	}

	for _, url := range opts.Args.Urls {
		client := *http.DefaultClient
		client.Timeout = opts.Timeout
		if !opts.Raw {
			fmt.Printf(`Dumping urls from: %s`+"\n", url)
		}

		urls, err := splitdumper.DumpWithWait(client, url, opts.Limit, func() {
			time.Sleep(opts.Wait)
		})
		if err != nil {
			fmt.Printf(`Error dumping domain on "%s": %v`, url, err)
			os.Exit(2)
			return
		}
		if !opts.Raw {
			fmt.Printf("Found %d destinations:\n", len(urls))
		}
		sort.Strings(urls)
		for _, dest := range urls {
			fmt.Println(dest)
		}
	}
}
