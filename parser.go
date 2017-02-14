package main

import (
	"io"
	"strings"
	"encoding/xml"
	"regexp"
)


func FindBody(content io.Reader) string {

	decoder := xml.NewDecoder(content)

	for {
		
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch el := t.(type) {
		case xml.StartElement:
			if el.Name.Local == "div" {

				for _, att := range el.Attr {
					if att.Name.Local == "id" && att.Value == "bodyContent" {
						query := struct {
							Text string `xml:",innerxml"`
						}{}
						decoder.DecodeElement(&query, &el)

						return query.Text
					}
				}
			}
		}
	}

	return ""
}

func FindLinks(content io.Reader) []string {

	decoder := xml.NewDecoder(content)
	links := make([]string, 0)
	linkFilter, _ := regexp.Compile("^/wiki/[a-zA-Z]+:.*$")
	linkValidation, _ := regexp.Compile("^/wiki/[a-zA-Z0-9_()*-]+$")

	for {
		
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch el := t.(type) {
		case xml.StartElement:
			if el.Name.Local == "a" {
				
				for _, att := range el.Attr {
					if att.Name.Local == "href" {
						link := att.Value
						if !linkFilter.MatchString(link) && linkValidation.MatchString(link) {
							links = append(links, link)
						}
					}
				}
			}
		}
	}

	return links
}

func Parse(content io.Reader, scorer Scorer) (float64, []string) {

	body := FindBody(content)

	scoreChan := make(chan float64, 1)
	go func (content string, result chan float64, scorer Scorer) {
		result <- scorer.GetScore(strings.NewReader(content))
	}(body, scoreChan, scorer)

	links := FindLinks(strings.NewReader(body))

	return <-scoreChan, links
}
