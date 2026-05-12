package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func filterURL(rawURL string, keysToRemove []string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", err
	}

	params := u.Query()
	for _, key := range keysToRemove {
		params.Del(key)
	}

	u.RawQuery = params.Encode()
	return u.String(), nil
}

func main() {
	// Define flags
	removePtr := flag.String("r", "", "Comma-separated list of parameters to remove")
	urlPtr := flag.String("u", "", "The single URL to process")
	flag.Parse()

	if *removePtr == "" {
		fmt.Println("Usage: paramcut -r itm_medium,itm_source [-u url]")
		os.Exit(1)
	}

	keys := strings.Split(*removePtr, ",")

	// Priority 1: Process a single URL via flag
	if *urlPtr != "" {
		result, err := filterURL(*urlPtr, keys)
		if err == nil {
			fmt.Println(result)
		}
		return
	}

	// Priority 2: Process via Stdin (piping)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			result, err := filterURL(scanner.Text(), keys)
			if err == nil && result != "" {
				fmt.Println(result)
			}
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
}