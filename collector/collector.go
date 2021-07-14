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
	IsInternal bool   `json:"is_internal_link"`
}

type HeadingCount struct {
	H1Count int `json:"h1_count"`
	H2Count int `json:"h2_count"`
	H3Count int `json:"h3_count"`
	H4Count int `json:"h4_count"`
	H5Count int `json:"h5_count"`
	H6Count int `json:"h6_count"`
}

type AppData struct {
	HtmlVersion string `json:"html_version"`
	Title string `json:"title"`
	HeadingCount HeadingCount `json:"heading_count"`
	Links []Link `json:"links"`
	HasLogin bool `json:"has_login"`
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

func (appData *AppData) AddLink(l Link) {
	appData.Links = append(appData.Links, l)
}

func GetAppData(url string) *AppData {
	appData := &AppData{}
	appData.setPageInfo(url)
	//appData.setLinkList(url)
	return appData
}

func (appData *AppData) setPageInfo(url string)  {
	c := colly.NewCollector(
		colly.IgnoreRobotsTxt(),
		colly.Async(true),
	)
	c.SetRequestTimeout(time.Second * 10)
	c.AllowURLRevisit = false

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML("html" , func(e *colly.HTMLElement) {
		//setting up the title
		appData.Title = e.ChildText("title")

		// set header tags count
		e.ForEach("h1, h2, h3, h4, h5, h6", func(_ int, el *colly.HTMLElement) {
			switch el.Name {
			case "h1":
				appData.HeadingCount.H1Count += 1
			case "h2":
				appData.HeadingCount.H2Count += 1
			case "h3":
				appData.HeadingCount.H3Count += 1
			case "h4":
				appData.HeadingCount.H4Count += 1
			case "h5":
				appData.HeadingCount.H5Count += 1
			case "h6":
				appData.HeadingCount.H6Count += 1
			}
		})

		// has login
		e.ForEach("input", func(i int, el *colly.HTMLElement) {
			if el.Attr("type") == "password"{
				appData.HasLogin = true
			}
		})
	})


	c.Visit(url)
	c.Wait()

	log.Printf("Scraping finished")
}

func (appData *AppData) setLinkList(url string) {

	depth := 2
	threads := 4

	flag.Parse()

	c := colly.NewCollector(
			colly.Async(true),
			colly.MaxDepth(depth),
			colly.IgnoreRobotsTxt(),
			colly.URLFilters(
			regexp.MustCompile("https?://.+$"),
		),
	)
	c.SetRequestTimeout(time.Second * 10)
	c.AllowURLRevisit = false

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

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(response *colly.Response) {
		baseUrl, _ := url2.Parse(url)
		link := Link{
			Url:        response.Request.URL.String(),
			Status:     response.StatusCode,
			IsInternal: isInternalLink(response, baseUrl),
		}
		appData.AddLink(link)
	})

	c.Visit(url)
	c.Wait()
}

// here subdomains considered as external
func isInternalLink(response *colly.Response, baseUrl *url2.URL) bool {
	return baseUrl.Host == response.Request.URL.Host
}

func handleError(error error) {
	if error != nil {
		fmt.Println("Error:", error)
	}
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



