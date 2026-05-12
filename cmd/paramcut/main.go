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

func filterURL(rawURL string, regexPatterns []*regexp.Regexp) (string, error) {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", err
	}

	params := u.Query()
	// Iterate through all parameters in the URL
	for pName := range params {
		for _, re := range regexPatterns {
			if re.MatchString(pName) {
				params.Del(pName)
				break 
			}
		}
	}

	u.RawQuery = params.Encode()
	
	// If the query string is empty, remove the '?' from the URL
	result := u.String()
	return strings.TrimSuffix(result, "?"), nil
}

func main() {
	removePtr := flag.String("r", "", "Comma-separated list of Regex patterns")
	urlPtr := flag.String("u", "", "Single URL to process")
	flag.Parse()

	if *removePtr == "" {
		fmt.Println("paramcut: Remove URL parameters that match the Regex patterns")
		fmt.Println("Version: 0.0.4")
		fmt.Println("Usage: paramcut -r \"utm_.*\" [-u url]")
		os.Exit(1)
	}

	// Fix: Trim spaces after splitting by comma
	rawPatterns := strings.Split(*removePtr, ",")
	var regexPatterns []*regexp.Regexp
	for _, p := range rawPatterns {
		p = strings.TrimSpace(p) // Important: remove hidden spaces
		if p == "" {
			continue
		}
		re, err := regexp.Compile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid regex: %s\n", p)
			continue
		}
		regexPatterns = append(regexPatterns, re)
	}

	// Handle single URL via -u flag
	if *urlPtr != "" {
		result, _ := filterURL(*urlPtr, regexPatterns)
		fmt.Println(result)
		return
	}

	// Handle multiple URLs via Stdin (piping)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		result, err := filterURL(line, regexPatterns)
		if err == nil {
			fmt.Println(result)
		}
	}
}