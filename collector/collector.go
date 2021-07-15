package collector

import (
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"strings"
	"time"
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
	HtmlVersion  string       `json:"html_version"`
	Title        string       `json:"title"`
	HeadingCount HeadingCount `json:"heading_count"`
	Links        []Link       `json:"links"`
	HasLogin     bool         `json:"has_login"`
}

func (appData *AppData) AddLink(link Link) {
	appData.Links = append(appData.Links, link)
}

// aggregate all results into appData
func GetAppData(url string) *AppData {
	appData := &AppData{}
	appData.setPageInfo(url)
	appData.setLinkList(url)
	return appData
}

// this aggregate html-version,
//header count and is login page results into the response.
func (appData *AppData) setPageInfo(url string) {
	collector := colly.NewCollector(
		colly.IgnoreRobotsTxt(),
		colly.Async(true),
	)
	collector.SetRequestTimeout(time.Second * COLLY_TIMEOUT_DURATION)
	collector.AllowURLRevisit = false

	collector.OnHTML("html", func(element *colly.HTMLElement) {
		// setting up the title
		appData.Title = element.ChildText("title")
		// set header tags count
		appData.setHeaderCount(element)
		// has login
		appData.setIsLogin(element)
	})

	collector.OnResponse(func(response *colly.Response) {
		appData.HtmlVersion = setHTMLVersion(string(response.Body))
	})

	collector.Visit(url)
	collector.Wait()
}

// counting the header types in the element
func (appData *AppData) setHeaderCount(element *colly.HTMLElement) {
	element.ForEach("h1, h2, h3, h4, h5, h6", func(_ int, headerElement *colly.HTMLElement) {
		switch headerElement.Name {
		case H1:
			appData.HeadingCount.H1Count += 1
		case H2:
			appData.HeadingCount.H2Count += 1
		case H3:
			appData.HeadingCount.H3Count += 1
		case H4:
			appData.HeadingCount.H4Count += 1
		case H5:
			appData.HeadingCount.H5Count += 1
		case H6:
			appData.HeadingCount.H6Count += 1
		}
	})
}

// check whether the element has a login element by type="password"
func (appData *AppData) setIsLogin(element *colly.HTMLElement) {
	element.ForEach("input", func(i int, inputElement *colly.HTMLElement) {
		if inputElement.Attr("type") == "password" {
			appData.HasLogin = true
		}
	})
}

// retrieving links associated with the
//web page with each one's status
// this uses separate colly collector
//due to it's different behaviour of functionality compared to others
func (appData *AppData) setLinkList(baseUrlStr string) {
	collector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(LINK_LIST_COLLECTOR_DEPTH),
		colly.IgnoreRobotsTxt(),
		colly.URLFilters(
			regexp.MustCompile("https?://.+$"),
		),
	)
	collector.SetRequestTimeout(time.Second * COLLY_TIMEOUT_DURATION)
	collector.AllowURLRevisit = false
	collector.AllowURLRevisit = false

	if limitError := collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: LINK_LIST_COLLECTOR_THREAD_COUNT,
		RandomDelay: 1 * time.Second,
	}); limitError != nil {
		log.WithFields(log.Fields{
			"error": limitError,
		}).Error("Collector LimitRule error")
	}
	//// Find and visit all links
	collector.OnHTML("a[href]", func(element *colly.HTMLElement) {
		element.Request.Visit(element.Attr("href"))
	})

	collector.OnError(func(response *colly.Response, err error) {
		log.WithFields(log.Fields{
			"url":        response.Request.URL.String(),
			"status":     response.StatusCode,
			"error": err,
		}).Debug("Navigated link error")

		baseUrl, _ := url.Parse(baseUrlStr)

		link := Link{
			Url:        response.Request.URL.String(),
			Status:     response.StatusCode,
			IsInternal: isInternalLink(response.Request.URL, baseUrl),
		}
		appData.AddLink(link)
	})

	collector.OnResponse(func(response *colly.Response) {
		baseUrl, _ := url.Parse(baseUrlStr)
		link := Link{
			Url:        response.Request.URL.String(),
			Status:     response.StatusCode,
			IsInternal: isInternalLink(response.Request.URL, baseUrl),
		}
		appData.AddLink(link)
	})
	collector.Visit(baseUrlStr)
	collector.Wait()
}

// here subdomains considered as external
func isInternalLink(url *url.URL, baseUrl *url.URL) bool {
	return baseUrl.Host == url.Host
}

// identify html version using the doctype
func setHTMLVersion(html string) string {
	var version = UNKNOWN
	for doctype, matcher := range GetHtmlVersions() {
		match := strings.Contains(html, matcher)
		if match == true {
			version = doctype
			break
		}
	}
	return version
}
