package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// filterURL parses the URL and removes query parameters that match any of the provided regex patterns.
func filterURL(rawURL string, regexPatterns []*regexp.Regexp) (string, error) {
	// Parse the raw string into a URL structure
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", err
	}

	// Extract query parameters
	params := u.Query()
	for pName := range params {
		for _, re := range regexPatterns {
			// If parameter name matches the regex, delete it
			if re.MatchString(pName) {
				params.Del(pName)
				break // Move to the next parameter once a match is found
			}
		}
	}

	// Re-encode the remaining parameters back into the URL
	u.RawQuery = params.Encode()
	return u.String(), nil
}

func main() {
	// Define CLI flags
	removePtr := flag.String("r", "", "Comma-separated list of Regex patterns to remove")
	urlPtr := flag.String("u", "", "The target URL to process")
	flag.Parse()

	// Validation: Ensure the removal pattern is provided
	if *removePtr == "" {
		fmt.Println("Usage: paramcut -r \"^utm_.*|.*_source$\" [-u url]")
		os.Exit(1)
	}

	// Pre-compile all regex patterns for efficiency
	rawPatterns := strings.Split(*removePtr, ",")
	var regexPatterns []*regexp.Regexp
	for _, p := range rawPatterns {
		re, err := regexp.Compile(strings.TrimSpace(p))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid regex pattern '%s': %v\n", p, err)
			continue
		}
		regexPatterns = append(regexPatterns, re)
	}

	// Priority 1: Process a single URL via flag
	if *urlPtr != "" {
		result, err := filterURL(*urlPtr, regexPatterns)
		if err == nil {
			fmt.Println(result)
		}
		return
	}

	// Priority 2: Process via Stdin (Piping support)
	// Check if data is being piped to the tool
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			result, err := filterURL(scanner.Text(), regexPatterns)
			if err == nil && result != "" {
				fmt.Println(result)
			}
		}
	} else {
		// If no URL flag and no piped input, show usage
		flag.Usage()
	}
}