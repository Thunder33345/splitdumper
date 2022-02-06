# Split Dumper

Dump all destinations of a split.to link

# Split Dumper Utility

You can install it with `go install github.com/Thunder33345/splitdumper/cmd/splitdumper`

## Split dumper Usage
`splitdumper http://url1 http://url2`
```
Usage:
splitdumper.exe [OPTIONS] [Urls...]

Application Options:
-t, --timeout:  Sets the client timeout (default: 5s)
-w, --wait:     Sets the wait time after crawling (default: 50ms)
-l, --limit:    How many times all url should be seen before stopping (default: 3)
-r, --raw       Outputs only the end result

Help Options:
-?             Show this help message
-h, --help      Show this help message

Arguments:
Urls:          The urls to dump
```
