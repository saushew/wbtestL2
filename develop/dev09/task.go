package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Wget .
type Wget struct {
	host  string
	depth int
	queue []string
	cache map[string]struct{}
	dir   string
}

var (
	depth      int
	websiteURL string
)

var (
	errorNoLink = errors.New("link is required")
)

func init() {
	flag.IntVar(&depth, "d", 1, "depth")

	flag.Parse()
	websiteURL = flag.Arg(0)
}

func main() {
	if websiteURL == "" {
		fmt.Printf("wget: %s\n", errorNoLink.Error())
		os.Exit(1)
	}

	w, err := NewWget(websiteURL, depth)
	if err != nil {
		fmt.Printf("wget: %s\n", err.Error())
		os.Exit(1)
	}

	err = w.Run()
	if err != nil {
		fmt.Printf("wget: %s\n", err.Error())
		os.Exit(1)
	}
}

// NewWget .
func NewWget(URL string, depth int) (*Wget, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	w := &Wget{
		host:  u.Host,
		depth: depth,
		queue: []string{URL},
		cache: map[string]struct{}{URL: {}},
		dir:   fmt.Sprintf("%s depth=%d", u.Host, depth),
	}

	return w, nil
}

// Run .
func (w *Wget) Run() error {
	if _, err := os.Stat(w.dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(w.dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	linkCntr := 1

	for len(w.queue) > 0 && w.depth > 0 {
		URL := w.queue[0]
		w.queue = w.queue[1:]

		w.processLink(URL)

		linkCntr--
		if linkCntr == 0 {
			linkCntr = len(w.queue)
			w.depth--
		}
	}

	return nil
}

func (w *Wget) processLink(URL string) error {
	r, err := http.Get(URL)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()

	saveHTML(data, w.dir, URL)

	b, err := html.Parse(bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	links := findLinks(b)

	ptr := 0
	for i, link := range links {
		_url, err := url.Parse(link)
		if err != nil {
			continue
		}

		if _url.Host != w.host {
			continue
		}

		if _, ok := w.cache[link]; !ok {
			w.cache[link] = struct{}{}
			links[i], links[ptr] = links[ptr], links[i]
			ptr++
		}
	}

	w.queue = append(w.queue, links[:ptr]...)
	return nil
}

func findLinks(n *html.Node) []string {
	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	return links
}

func saveHTML(data []byte, dir string, URL string) {
	filename := strings.ReplaceAll(URL, "/", "|") + ".html"

	p := path.Join(dir, filename)

	f, err := os.Create(p)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(data)
}
