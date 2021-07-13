package collector

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	url2 "net/url"
	"regexp"
	"strings"
	"time"
)

var (
	appData AppData
)

type Link struct {
	Url        string `json:"url"`
	Status     int    `json:"status"`
	IsInternal bool   `json:"isInternal"`
}

type AppData struct {
	Title string `json:"title"`
	Links []Link `json:"links"`
}

var docTypes = make(map[string]string)

func init() {
	docTypes["HTML 4.01 Strict"] = `"-//W3C//DTD HTML 4.01//EN"`
	docTypes["HTML 4.01 Transitional"] = `"-//W3C//DTD HTML 4.01 Transitional//EN"`
	docTypes["HTML 4.01 Frameset"] = `"-//W3C//DTD HTML 4.01 Frameset//EN"`
	docTypes["XHTML 1.0 Strict"] = `"-//W3C//DTD XHTML 1.0 Strict//EN"`
	docTypes["XHTML 1.0 Transitional"] = `"-//W3C//DTD XHTML 1.0 Transitional//EN"`
	docTypes["XHTML 1.0 Frameset"] = `"-//W3C//DTD XHTML 1.0 Frameset//EN"`
	docTypes["XHTML 1.1"] = `"-//W3C//DTD XHTML 1.1//EN"`
	docTypes["HTML 5"] = `<!DOCTYPE html>`
}

func (a *AppData) AddLink(l Link) {
	a.Links = append(a.Links, l)
}

func getAppData(url string) *AppData {

}

func GetLinkList(url string) *AppData {

	depth := 2
	threads := 4

	flag.Parse()

	c := colly.NewCollector(
			colly.Async(true),
			colly.MaxDepth(depth),
			colly.URLFilters(
			regexp.MustCompile("https?://.+$"),
		),
	)

	limitError := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: threads,
		RandomDelay: 1 * time.Second,
	})

	handleError(limitError)

	//// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("found title", e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(response *colly.Response) {
		baseUrl, _ := url2.Parse(url)
		link := Link{
			Url:        response.Request.URL.String(),
			Status:     response.StatusCode,
			IsInternal: isInternal(response, baseUrl),
		}
		appData.AddLink(link)
		log.Println("app data ", appData)
	})

	c.Visit(url)
	c.Wait()

	return &appData
}



func checkDoctype(html string) string {
	var version = "UNKNOWN"

	for doctype, matcher := range docTypes {
		match := strings.Contains(html, matcher)

		if match == true {
			version = doctype
		}
	}

	return version
}


func isInternal(response *colly.Response, baseUrl *url2.URL) bool {
	return baseUrl.Host == response.Request.URL.Host
}

func handleError(error error) {
	if error != nil {
		fmt.Println("Error:", error)
	}
}
