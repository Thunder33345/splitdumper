package main

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/thunder33345/splitdumper"
	"net/http"
	"os"
	"os/signal"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handleExit(func() {
		cancel()
	})

	client := *http.DefaultClient
	client.Timeout = opts.Timeout
	sleeper := func() {
		time.Sleep(opts.Wait)
	}
	for _, url := range opts.Args.Urls {
		if !opts.Raw {
			fmt.Printf(`Dumping urls from: %s`+"\n", url)
		}

		urls, err := splitdumper.Dump(url, opts.Limit, splitdumper.WithClient(client), splitdumper.WithWait(sleeper), splitdumper.WithContext(ctx))
		if err != nil {
			fmt.Printf(`Error dumping domain on "%s": %v`, url, err)
			fmt.Println()
		}
		if !opts.Raw {
			fmt.Printf("Found %d destinations:\n", len(urls))
		}
		sort.Strings(urls)
		for _, dest := range urls {
			fmt.Println(dest)
		}
		if err != nil {
			fmt.Printf("Partial results")
			os.Exit(2)
			return
		}
	}
}

func handleExit(cb func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cb()
		}
	}()
}
