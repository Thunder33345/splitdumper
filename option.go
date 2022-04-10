package splitdumper

import (
	"context"
	"net/http"
)

//config is an internal object that stores the configuration for the Dump function
//this is internal representation, should not be exported or exposed.
type config struct {
	// client the http client to use for the requests
	client *http.Client
	// url the url to dump
	url string
	// limit is the minimum times a destination should be seen
	limit int
	// context is for cancellation
	context context.Context
	// wait will be called everytime the job an url is obtained, can be used to pause between request
	wait func()
	// hook allows application to get realtime progress
	hook    func(url string, seen int)
	breaker Breaker
}

//Option is a function that can be used to configure the Dump function
type Option func(*config)

//WithClient configures a http client to use
func WithClient(client *http.Client) Option {
	return func(c *config) {
		c.client = client
	}
}

//WithContext configures a context to use
func WithContext(ctx context.Context) Option {
	return func(c *config) {
		c.context = ctx
	}
}

//WithWait configures a function to call everytime an url is obtained
func WithWait(wait func()) Option {
	return func(c *config) {
		if wait == nil {
			wait = func() {}
		}
		c.wait = wait
	}
}

//WithHook configures a function to call everytime an url is obtained
func WithHook(hook func(url string, seen int)) Option {
	return func(c *config) {
		if hook == nil {
			hook = func(url string, seen int) {}
		}
		c.hook = hook
	}
}
