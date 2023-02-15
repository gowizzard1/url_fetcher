package utils

import (
	"flag"
	"fmt"
	"os"
)

func ParseFlagsAndArgsToURL() []string {
	flag.Parse()

	urls := flag.Args()
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: url-sizes [options] <url>...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return urls
}
