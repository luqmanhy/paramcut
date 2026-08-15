# paramcut ✂️

`paramcut` is a fast, lightweight CLI tool written in Go designed to strip unwanted query parameters from URLs. It works by using regular expressions to match and remove the specified query parameters from the URLs. You can define the query parameters to remove in the `-r` flag, and the URLs to process in the `-u` flag.

# Installation
```bash
go install -v github.com/luqmanhy/paramcut/cmd/paramcut@latest
```

# Usage
```bash
# Remove all utm_.* query parameters from a single URL
paramcut -r "utm_.*" -u "https://example.com/?utm_source=google&utm_medium=pc&utm_campaign=1234567890&param1=value1"

# Output: https://example.com/?param1=value1

# Process multiple URLs from a file
cat urls.txt | paramcut -r "utm_.*"
```