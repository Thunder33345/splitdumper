# Split Dumper

Dump all destinations of a split.to link

# Split Dumper Utility

You can install it with `go install github.com/Thunder33345/splitdumper/cmd/splitdumper@latest`

## Split Dumper Usage
`splitdumper http://url1 http://url2`

```
Usage:
  splitdumper.exe [OPTIONS] [Urls...]

Application Options:
  -t, --timeout= Sets the client timeout (default: 5s)
  -w, --wait=    Sets the wait time after crawling (default: 50ms)
  -l, --limit=   How many times all url should be seen before stopping (default: 3)
  -f, --full     Follow the redirect for final destination(will access out of the URLs host)
  -r, --raw      Raw format only the end result, and Error'
      --text     (Default)Text format meant for human consumption'
  -j, --json     JSON format that can be parsed by other tools'

Help Options:
  -h, --help     Show this help message

Arguments:
  Urls:          The urls to dump (required)
```

# Split Dumper TUI

The TUI version can be installed with `go install github.com/Thunder33345/splitdumper/cmd/splitdumpertui@latest`

There's not much difference except a terminal user interface is provided