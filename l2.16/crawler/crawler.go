package crawler

import (
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"wget/downloader"
	parser "wget/pageparser"

	"golang.org/x/net/html"
)

type Crawler struct {
	downloader *downloader.Downloader
	parser     *parser.PageParser
	maxDepth   int
	outputDir  string

	mu         sync.Mutex
	downloaded map[string]struct{}
}

func NewCrawler(startURL string, maxDepth int, outputDir string) (*Crawler, error) {
	pageParser, err := parser.NewPageParser(startURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		downloader: downloader.NewDownloader(10 * time.Second),
		parser:     pageParser,
		maxDepth:   maxDepth,
		outputDir:  outputDir,
		downloaded: make(map[string]struct{}),
	}, nil
}

func (c *Crawler) Start(startURL string) error {
	return c.crawl(startURL, 0, true)
}

func (c *Crawler) crawl(targetURL string, currentDepth int, followPageLinks bool) error {
	if currentDepth > c.maxDepth {
		return nil
	}

	parsedURL, err := c.parser.ResolveURL(nil, targetURL)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", targetURL, err)
	}

	if !c.parser.IsSameDomain(parsedURL) {
		return nil
	}

	parsedURL.Fragment = ""
	if !c.markDownloaded(parsedURL.String()) {
		return nil
	}

	resp, err := c.downloader.Download(parsedURL.String())
	if err != nil {
		return fmt.Errorf("download %s: %w", parsedURL.String(), err)
	}
	defer resp.Body.Close()

	contentType := responseMediaType(resp.Header.Get("Content-Type"))
	isHTML := contentType == "text/html"
	isCSS := contentType == "text/css"
	localPath := c.parser.URLToLocalPath(parsedURL, isHTML)
	fullPath := filepath.Join(c.outputDir, localPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(fullPath), err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", fullPath, err)
	}
	defer file.Close()

	switch {
	case isHTML:
		if err := c.saveHTMLPage(file, resp.Body, parsedURL, localPath, currentDepth, followPageLinks); err != nil {
			return err
		}
	case isCSS:
		if err := c.saveCSSFile(file, resp.Body, parsedURL, localPath, currentDepth); err != nil {
			return err
		}
	default:
		if _, err := io.Copy(file, resp.Body); err != nil {
			return fmt.Errorf("write %s: %w", fullPath, err)
		}
	}

	return nil
}

func (c *Crawler) saveHTMLPage(file *os.File, body io.Reader, pageURL *url.URL, localPath string, currentDepth int, followPageLinks bool) error {
	doc, err := html.Parse(body)
	if err != nil {
		return fmt.Errorf("parse HTML %s: %w", pageURL.String(), err)
	}

	parseResult := c.parser.RewriteHTML(doc, pageURL, localPath)
	if err := html.Render(file, doc); err != nil {
		return fmt.Errorf("write HTML %s: %w", localPath, err)
	}

	for _, resourceURL := range parseResult.Resources {
		if err := c.crawl(resourceURL, currentDepth, false); err != nil {
			fmt.Fprintf(os.Stderr, "resource error: %v\n", err)
		}
	}

	if followPageLinks && currentDepth < c.maxDepth {
		for _, pageLink := range parseResult.Pages {
			if err := c.crawl(pageLink, currentDepth+1, true); err != nil {
				fmt.Fprintf(os.Stderr, "page error: %v\n", err)
			}
		}
	}

	return nil
}

func (c *Crawler) saveCSSFile(file *os.File, body io.Reader, cssURL *url.URL, localPath string, currentDepth int) error {
	content, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read CSS %s: %w", cssURL.String(), err)
	}

	rewrittenCSS, assets := c.parser.RewriteCSS(string(content), cssURL, localPath)
	if _, err := io.WriteString(file, rewrittenCSS); err != nil {
		return fmt.Errorf("write CSS %s: %w", localPath, err)
	}

	for _, assetURL := range assets {
		if err := c.crawl(assetURL, currentDepth, false); err != nil {
			fmt.Fprintf(os.Stderr, "css asset error: %v\n", err)
		}
	}

	return nil
}

func (c *Crawler) markDownloaded(targetURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.downloaded[targetURL]; ok {
		return false
	}

	c.downloaded[targetURL] = struct{}{}
	return true
}

func responseMediaType(contentType string) string {
	if contentType == "" {
		return ""
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return strings.ToLower(contentType)
	}

	return strings.ToLower(mediaType)
}
