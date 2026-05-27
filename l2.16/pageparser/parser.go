package parser

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var (
	cssURLPattern    = regexp.MustCompile(`url\(\s*(['"]?)([^'")]+)\1\s*\)`)
	cssImportPattern = regexp.MustCompile(`@import\s+(?:url\(\s*)?(['"]?)([^'")\s]+)\1\s*\)?`)
)

type ParseResult struct {
	Pages     []string
	Resources []string
}

type PageParser struct {
	baseDomain string
}

func NewPageParser(startURL string) (*PageParser, error) {
	u, err := url.Parse(startURL)
	if err != nil {
		return nil, err
	}

	return &PageParser{baseDomain: normalizedHost(u.Host)}, nil
}

func (p *PageParser) ResolveURL(base *url.URL, raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty URL")
	}

	if strings.HasPrefix(raw, "#") {
		return nil, fmt.Errorf("fragment-only URL")
	}

	lowerRaw := strings.ToLower(raw)
	if strings.HasPrefix(lowerRaw, "mailto:") ||
		strings.HasPrefix(lowerRaw, "javascript:") ||
		strings.HasPrefix(lowerRaw, "tel:") ||
		strings.HasPrefix(lowerRaw, "data:") {
		return nil, fmt.Errorf("unsupported URL scheme")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if base != nil {
		parsed = base.ResolveReference(parsed)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme %q", parsed.Scheme)
	}

	if parsed.Host == "" {
		return nil, fmt.Errorf("missing host")
	}

	parsed.Fragment = strings.TrimSpace(parsed.Fragment)
	return parsed, nil
}

func (p *PageParser) IsSameDomain(targetURL *url.URL) bool {
	return normalizedHost(targetURL.Host) == p.baseDomain
}

func (p *PageParser) URLToLocalPath(targetURL *url.URL, asHTML bool) string {
	clean := cloneURL(targetURL)
	clean.Fragment = ""

	localPath := clean.EscapedPath()
	if localPath == "" {
		localPath = "/"
	}

	if asHTML {
		localPath = htmlPath(localPath)
	} else {
		localPath = assetPath(localPath)
	}

	localPath = addQuerySuffix(localPath, clean.RawQuery)
	localPath = strings.TrimPrefix(path.Clean(localPath), "/")
	if localPath == "." || localPath == "" {
		if asHTML {
			localPath = "index.html"
		} else {
			localPath = "index"
		}
	}

	return filepath.Join(sanitizeHost(targetURL.Host), filepath.FromSlash(localPath))
}

func (p *PageParser) RewriteHTML(doc *html.Node, pageURL *url.URL, currentLocalPath string) ParseResult {
	result := ParseResult{}
	seenPages := make(map[string]struct{})
	seenResources := make(map[string]struct{})

	var visit func(node *html.Node)
	visit = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for i := range node.Attr {
				attr := &node.Attr[i]
				linkKind, ok := classifyHTMLLink(node.Data, attr.Key)
				if !ok {
					continue
				}

				resolved, err := p.ResolveURL(pageURL, attr.Val)
				if err != nil || !p.IsSameDomain(resolved) {
					continue
				}

				attr.Val = p.relativeLocalReference(currentLocalPath, resolved, linkKind == linkPage)

				canonical := cloneURL(resolved)
				canonical.Fragment = ""
				if linkKind == linkPage {
					appendUnique(&result.Pages, seenPages, canonical.String())
				} else {
					appendUnique(&result.Resources, seenResources, canonical.String())
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}

	visit(doc)
	return result
}

func (p *PageParser) RewriteCSS(content string, cssURL *url.URL, currentLocalPath string) (string, []string) {
	discovered := make([]string, 0)
	seen := make(map[string]struct{})

	rewrite := func(match string, expression *regexp.Regexp, isImport bool) string {
		parts := expression.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}

		resolved, err := p.ResolveURL(cssURL, parts[2])
		if err != nil || !p.IsSameDomain(resolved) {
			return match
		}

		canonical := cloneURL(resolved)
		canonical.Fragment = ""
		appendUnique(&discovered, seen, canonical.String())

		localRef := p.relativeLocalReference(currentLocalPath, resolved, false)
		if isImport {
			return fmt.Sprintf(`@import "%s"`, localRef)
		}
		return fmt.Sprintf(`url("%s")`, localRef)
	}

	content = cssImportPattern.ReplaceAllStringFunc(content, func(match string) string {
		return rewrite(match, cssImportPattern, true)
	})
	content = cssURLPattern.ReplaceAllStringFunc(content, func(match string) string {
		return rewrite(match, cssURLPattern, false)
	})

	return content, discovered
}

type htmlLinkKind int

const (
	linkPage htmlLinkKind = iota
	linkResource
)

func classifyHTMLLink(tagName, attrKey string) (htmlLinkKind, bool) {
	switch tagName {
	case "a", "iframe":
		if attrKey == "href" || attrKey == "src" {
			return linkPage, true
		}
	case "img", "script", "source", "audio":
		if attrKey == "src" {
			return linkResource, true
		}
	case "video":
		if attrKey == "src" || attrKey == "poster" {
			return linkResource, true
		}
	case "link":
		if attrKey == "href" {
			return linkResource, true
		}
	}

	return 0, false
}

func (p *PageParser) relativeLocalReference(fromLocalPath string, targetURL *url.URL, asHTML bool) string {
	clean := cloneURL(targetURL)
	fragment := clean.Fragment
	clean.Fragment = ""

	targetLocalPath := p.URLToLocalPath(clean, asHTML)
	relativePath, err := filepath.Rel(filepath.Dir(fromLocalPath), targetLocalPath)
	if err != nil {
		relativePath = targetLocalPath
	}

	relativePath = filepath.ToSlash(relativePath)
	if relativePath == "." {
		relativePath = filepath.Base(targetLocalPath)
	}

	if fragment != "" {
		relativePath += "#" + fragment
	}

	return relativePath
}

func appendUnique(target *[]string, seen map[string]struct{}, value string) {
	if _, ok := seen[value]; ok {
		return
	}

	seen[value] = struct{}{}
	*target = append(*target, value)
}

func normalizedHost(host string) string {
	return strings.TrimPrefix(strings.ToLower(host), "www.")
}

func sanitizeHost(host string) string {
	replacer := strings.NewReplacer(":", "_")
	return replacer.Replace(host)
}

func htmlPath(rawPath string) string {
	switch {
	case rawPath == "" || rawPath == "/":
		return "/index.html"
	case strings.HasSuffix(rawPath, "/"):
		return rawPath + "index.html"
	case path.Ext(rawPath) == "":
		return rawPath + "/index.html"
	default:
		return rawPath
	}
}

func assetPath(rawPath string) string {
	switch {
	case rawPath == "" || rawPath == "/":
		return "/index"
	case strings.HasSuffix(rawPath, "/"):
		return rawPath + "index"
	default:
		return rawPath
	}
}

func addQuerySuffix(localPath, rawQuery string) string {
	if rawQuery == "" {
		return localPath
	}

	queryHash := sha1.Sum([]byte(rawQuery))
	suffix := "__q_" + hex.EncodeToString(queryHash[:6])
	ext := path.Ext(localPath)
	if ext == "" {
		return localPath + suffix
	}

	base := strings.TrimSuffix(localPath, ext)
	return base + suffix + ext
}

func cloneURL(original *url.URL) *url.URL {
	copy := *original
	return &copy
}
