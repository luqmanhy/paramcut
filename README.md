# paramcut ✂️

`paramcut` is a fast, lightweight CLI tool written in Go designed to strip unwanted query parameters from URLs. It works by using regular expressions to match and remove the specified query parameters from the URLs. You can define the query parameters to remove in the `-r` flag, and the URLs to process in the `-u` flag.

# Usage
```bash
paramcut -r "utm_.*" -u url

cat urls.txt | paramcut -r "utm_.*"
```