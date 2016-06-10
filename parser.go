package main

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

/**
 * We keep a node (wikipedia link) only if it does not contains ':'
 * (to exclude pages on categories, or files), nor '_(disambiguation)'
 * which are not real entities.
 * More rules could probably be added.
 */
func isValidNode(link string) bool {
	if strings.HasPrefix(link, "/wiki/") && !strings.HasSuffix(link, "_(disambiguation)") &&  ! strings.Contains(link, ":") {
		return true
	}
	return false
}

/**
 * Parse the given html page, fetched from en.wikipedia.com
 * Basically only retrieve href elements from <a/> tags and
 * filter them (based on handmade rules) to only keep pages
 * about entities (no file, menu, ...).
 */
func parse(content io.ReadCloser) []string {

	defer content.Close()

	results := make([]string, 0)

	parser := html.NewTokenizer(content)
	inContent := false
	divCount := 0

	for {
		evt := parser.Next()
		switch {
		case evt == html.ErrorToken:
			break
		case evt == html.StartTagToken:
			t := parser.Token()
			if t.Data == "div" {
				
				if inContent {
					divCount += 1
				} else {
					for _, attr := range t.Attr {
	    				if attr.Key == "id" && attr.Val == "bodyContent" {
	        				inContent = true
	        				break
	    				}
					}
				}
			
			} else if t.Data == "a" && inContent {
				for _, attr := range t.Attr {
    				if attr.Key == "href" && isValidNode(attr.Val) {
    					results = append(results, attr.Val)
        				break
    				}
				}
			}
			break
		case evt == html.EndTagToken:
			t := parser.Token()
			if t.Data == "div" && inContent {
				if divCount == 0 {
					return results
				} else {
					divCount -= 1
				}
			}
			break
		}
	}

	return results
}