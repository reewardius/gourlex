package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	var filePath string
	var cookie string
	var customHeader string
	var proxyFlag string
	var urlOnly bool
	var pathOnly bool
	var silentMode bool
	flag.StringVar(&filePath, "f", "", "Specify file containing URLs")
	flag.StringVar(&cookie, "c", "", "Specify cookies")
	flag.StringVar(&customHeader, "r", "", "Specify headers")
	flag.StringVar(&proxyFlag, "p", "", "Specify the proxy URL")
	flag.BoolVar(&urlOnly, "uO", false, "Extract only URLs")
	flag.BoolVar(&pathOnly, "pO", false, "Extract only paths")
	flag.BoolVar(&silentMode, "s", false, "Silent mode")
	helpFlag := flag.Bool("h", false, "Display help")
	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	if !silentMode {
		printBanner()
	}

	if filePath == "" {
		fmt.Println("Error: No input file specified.")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	client := createHTTPClient(proxyFlag, silentMode)
	for scanner.Scan() {
		url := scanner.Text()
		processURL(url, client, cookie, customHeader, urlOnly, pathOnly, silentMode)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from file: %v\n", err)
	}
}

func processURL(url string, client *http.Client, cookie, customHeader string, urlOnly, pathOnly, silentMode bool) {
	if url == "" {
		return
	}
	validUrl, err := validateUrl(url)
	if err != nil {
		fmt.Printf("Error validating URL: %v\n", err)
		return
	}
	req, err := http.NewRequest("GET", validUrl, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	setupRequestHeaders(req, cookie, customHeader)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	urls, paths, err := extractURLsAndPaths(resp)
	if err != nil {
		fmt.Printf("Error extracting URLs and paths: %v\n", err)
		return
	}
	printResults(urls, paths, urlOnly, pathOnly, silentMode)
}

func createHTTPClient(proxyFlag string, silentMode bool) *http.Client {
	if proxyFlag != "" {
		proxyURL, err := url.Parse(proxyFlag)
		if err != nil {
			fmt.Printf("Error parsing proxy URL: %v\n", err)
			return &http.Client{}
		}
		if !silentMode {
			fmt.Printf("Using proxy: %s\n", proxyFlag)
		}
		return &http.Client{
			Transport: &http.Transport{
				Proxy:           http.ProxyURL(proxyURL),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	return &http.Client{}
}

func setupRequestHeaders(req *http.Request, cookie, customHeader string) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	if customHeader != "" {
		parts := strings.SplitN(customHeader, ":", 2)
		if len(parts) == 2 {
			req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
}

func extractURLsAndPaths(resp *http.Response) ([]string, []string, error) {
	tokenizer := html.NewTokenizer(resp.Body)
	var urls, paths []string
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return urls, paths, nil
		}
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			for _, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					if u, err := url.Parse(attr.Val); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
						urls = append(urls, u.String())
					} else {
						paths = append(paths, attr.Val)
					}
				}
			}
		}
	}
}

func printResults(urls, paths []string, urlOnly, pathOnly, silentMode bool) {
	if !silentMode {
		fmt.Printf("Extracted URLs from page:\n")
	}
	if !pathOnly {
		for _, url := range urls {
			fmt.Println(url)
		}
	}
	if !urlOnly {
		fmt.Println("\nPaths found on the page:")
		for _, path := range paths {
			fmt.Println(path)
		}
	}
}

func validateUrl(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	return u.String(), nil
}

func printHelp() {
	fmt.Println("Usage: ./gourlex -f <file_path> [options]")
	fmt.Println("Options:")
	fmt.Println("  -f string   Specify file containing URLs")
	fmt.Println("  -c string   Specify cookies")
	fmt.Println("  -r string   Specify custom headers")
	fmt.Println("  -p string   Specify the proxy URL")
	fmt.Println("  -uO         Extract only URLs")
	fmt.Println("  -pO         Extract only paths")
	fmt.Println("  -s          Silent mode (suppress output)")
	fmt.Println("  -h          Display this help and exit")
}

func printBanner() {
	fmt.Println("Gourlex - WebPage Urls Extractor Tool\n")
}
