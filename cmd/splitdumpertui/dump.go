package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thunder33345/splitdumper"
	"net/http"
	"time"
)

func dump(ch chan tea.Msg, url string, limit int, timeout, wait time.Duration) func() {
	c := *http.DefaultClient
	c.Timeout = timeout
	ctx, cancel := context.WithCancel(context.Background())
	waitFn := func() {
		time.Sleep(wait)
	}
	hook := func(url string, seen int) {
		ch <- hookMsg{
			url:  url,
			seen: seen,
		}
	}
	go func() {
		defer close(ch)
		urls, err := splitdumper.Dump(url, limit,
			splitdumper.WithClient(&c), splitdumper.WithContext(ctx), splitdumper.WithWait(waitFn), splitdumper.WithHook(hook))
		ch <- resultMsg{
			urls: urls,
			err:  err,
		}
	}()
	return cancel
}

type resultMsg struct {
	urls []string
	err  error
}
type hookMsg struct {
	url  string
	seen int
}

func waitForActivity(ch chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
