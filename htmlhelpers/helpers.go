package htmlhelpers

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func node(tag string, doc *html.Node) (*html.Node, error) {
	var n *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tag {
			n = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if n != nil {
		return n, nil
	}
	return nil, fmt.Errorf("missing %s in the node tree", tag)
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func extractLinks(content string) ([]string, error) {
	var (
		err       error
		links     []string = make([]string, 0)
		matches   [][]string
		findLinks = regexp.MustCompile("<a.*?href=\"(.*?)\"")
	)
	// Retrieve all anchor tag URLs from string
	matches = findLinks.FindAllStringSubmatch(content, -1)
	for _, val := range matches {
		var linkUrl *url.URL

		// Parse the anchr tag URL
		if linkUrl, err = url.Parse(val[1]); err != nil {
			return links, err
		}
		links = append(links, linkUrl.String())

	}
	return links, nil
}

func getHTMLSubset(tag string, content string) (string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", err
	}
	bn, err := node(tag, doc)
	if err != nil {
		return "", err
	}
	return renderNode(bn), nil
}

func GetLinks(content string) ([]string, error) {
	linkhtml, err := getHTMLSubset("ul", content)
	if err != nil {
		return []string{}, err
	}
	return extractLinks(linkhtml)

}
