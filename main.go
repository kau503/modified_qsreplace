package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

func main() {
	var appendMode bool
	flag.BoolVar(&appendMode, "a", false, "Append the value instead of replacing it")

	var ignorePath bool
	flag.BoolVar(&ignorePath, "ignore-path", false, "Ignore the path when considering what constitutes a duplicate")
	flag.Parse()

	seen := make(map[string]bool)

	// read URLs on stdin, then replace the values in the query string
	// with some user-provided value
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse url %s [%s]\n", sc.Text(), err)
			continue
		}

		// Go's maps aren't ordered, but we want to use all the param names
		// as part of the key to output only unique requests. To do that, put
		// them into a slice and then sort it.
		pp := make([]string, 0)
		for p, _ := range u.Query() {
			pp = append(pp, p)
		}
		sort.Strings(pp)

		key := fmt.Sprintf("%s%s?%s", u.Hostname(), u.EscapedPath(), strings.Join(pp, "&"))
		if ignorePath {
			key = fmt.Sprintf("%s?%s", u.Hostname(), strings.Join(pp, "&"))
		}

		// Only output each host + path + params combination once
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = true

		// Build query string manually without encoding
		var newQuery []string
		for param, vv := range u.Query() {
			if appendMode {
				newQuery = append(newQuery, fmt.Sprintf("%s=%s%s", param, vv[0], flag.Arg(0)))
			} else {
				newQuery = append(newQuery, fmt.Sprintf("%s=%s", param, flag.Arg(0)))
			}
		}
		u.RawQuery = strings.Join(newQuery, "&")

		fmt.Printf("%s\n", u)
	}
}
