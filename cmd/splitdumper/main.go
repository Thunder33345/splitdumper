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
	Full    bool          `short:"f" long:"full" description:"Follow the redirect for final destination(will access out of the URLs host)"`
	Raw     bool          `short:"r" long:"raw" description:"Raw format only the end result, and Error'"`
	Text    bool          `long:"text" description:"(Default)Text format meant for human consumption'"`
	JSON    bool          `short:"j" long:"json" description:"JSON format that can be parsed by other tools'"`
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

	ft, isOk := getFormatter()
	if !isOk {
		fmt.Printf("Cannot select multiple formatter")
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

	options := []splitdumper.Option{
		splitdumper.WithClient(&client), splitdumper.WithWait(sleeper), splitdumper.WithContext(ctx),
	}
	if opts.Full {
		options = append(options, splitdumper.WithFullTrace())
	}

	var all []Result
	var hasError bool

	for _, url := range opts.Args.Urls {
		output(ft.Start(url))
		urls, errD := splitdumper.Dump(url, opts.Limit, options...)

		sort.Strings(urls)

		output(ft.Result(url, urls, errD))
		all = append(all, Result{Url: url, Destinations: urls, Error: errD})
		if errD != nil {
			hasError = true
			if errD == ctx.Err() {
				break
			}
		}
	}
	output(ft.Complete(all))
	if hasError {
		os.Exit(2)
	}
}

var hasNewline = true

func output(str string) {
	if str == "" {
		return
	}
	if !hasNewline {
		fmt.Print("\n")
	} else {
		hasNewline = false
	}
	fmt.Print(str)
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

func getFormatter() (formatter, bool) {
	var ft formatter
	if opts.Raw {
		if ft != nil {
			return ft, false
		}
		ft = rawFormat{}
	}
	if opts.Text {
		if ft != nil {
			return ft, false
		}
		ft = textFormat{}
	}
	if opts.JSON {
		if ft != nil {
			return ft, false
		}
		ft = jsonFormat{}
	}

	if ft == nil {
		ft = textFormat{}
	}
	return ft, true
}
