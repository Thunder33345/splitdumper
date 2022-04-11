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
	// fullTrace if true, dump will follow and visit external sites to get the final destination
	// false by default, which only accesses the url to get the destination via "location" header, without visiting it
	fullTrace bool
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
//Depending on trace mode the client's CheckRedirect function could be overwritten
//If the client has to be overwritten, it will be dereference-ed first
func WithClient(client *http.Client) Option {
	return func(c *config) {
		c.client = client
	}
}

//WithFullTrace configures the Dump function to visit external sites to get the final destination
//by default it will only send HEAD request to the url and get the location via "location" header
func WithFullTrace() Option {
	return func(c *config) {
		c.fullTrace = true
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
