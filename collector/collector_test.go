package collector

import (
	"net/url"
	"reflect"
	"testing"
)

var sampleHTML5 = "<!DOCTYPE html>\n<head>\n<title>A Sample HTML Document (Test File)</title>\n<meta charset=" +
	"\"utf-8\">\n<meta name=\"description\" content=\"A blank HTML document for testing purposes.\">\n<meta name=" +
	"\"author\" content=\"Six Revisions\">\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1" +
	"\">\n<link rel=\"icon\" href=\"http://sixrevisions.com/favicon.ico\" type=\"image/x-icon\" />\n</head>\n<body>\n" +
	"  \n<h1>A Sample HTML Document (Test File)</h1>\n<p>A blank HTML document for testing purposes.</p>\n<p><a href=" +
	"\"../html5download-demo.html\">Go back to the demo</a></p>\n<p><a href=\"http://sixrevisions.com/html5/download-" +
	"attribute/\">Read the HTML5 download attribute guide</a></p>\n\n</body>\n</html>"

var sampleHTML4 = "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\"\n        \"http://www.w3.org/TR/html4/strict.dtd" +
	"\">\n<HTML>\n  <HEAD>\n    <TITLE>The document title</TITLE>\n  </HEAD>\n  <BODY>\n    <H1>Main heading</H1>\n   " +
	" <P>A paragraph.</P>\n    <P>Another paragraph.</P>\n    <UL>\n      <LI>A list item.</LI>\n      <LI>Another list" +
	" item.</LI>\n    </UL>\n  </BODY>\n</HTML>"

var complexHtmlHeaderCount = HeadingCount{
	H1Count: 0,
	H2Count: 2,
	H3Count: 8,
	H4Count: 4,
	H5Count: 0,
	H6Count: 0,
}

var complexHTMLLinks = []Link{
	{
		Url:        "http://eshankavishka.tech/",
		Status:     200,
		IsInternal: true,
	},
	{
		Url:        "https://twitter.com/Eshxnk",
		Status:     200,
		IsInternal: false,
	},
	{
		Url:        "https://github.com/eshan10x/Gym-management-system-with-javaFx",
		Status:     200,
		IsInternal: false,
	},
	{
		Url:        "https://github.com/Induw/UAsk",
		Status:     200,
		IsInternal: false,
	},
	{
		Url:        "https://github.com/eshan10x/Premipremier-League-Manager-System",
		Status:     200,
		IsInternal: false,
	},
	{
		Url:        "https://github.com/eshan10x/MY_MOVIE_APPLICATION_ANDROID",
		Status:     200,
		IsInternal: false,
	},
	{
		Url:        "https://drive.google.com/file/d/1R3j2AjX_UT2NpYz4wbFrBV59XAeZwv29/view?usp=sharing",
		Status:     200,
		IsInternal: false,
	},
	{
		Url:        "https://github.com/eshan10x/Simple-Finance-Calculator-With-JavaFx",
		Status:     200,
		IsInternal: false,
	},
	{
		Url:        "https://github.com/eshan10x",
		Status:     200,
		IsInternal: false,
	},
}

func Test_isInternalLink1(t *testing.T) {

	baseUrl, _ := url.Parse("https://www.home24.de/sale-ueberblick/")
	internalUrl, _ := url.Parse("https://www.home24.de/serviceversprechen/")
	externalUrl, _ := url.Parse("https://www.facebook.com/home24.de")
	subDomain, _ := url.Parse("https://www.sub.home24.de/sale-ueberblick/")

	type args struct {
		url     *url.URL
		baseUrl *url.URL
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"base_url", args{url: baseUrl, baseUrl: baseUrl}, true},
		{"internal_url", args{url: internalUrl, baseUrl: baseUrl}, true},
		{"external_url", args{url: externalUrl, baseUrl: baseUrl}, false},
		{"subdomain", args{url: subDomain, baseUrl: baseUrl}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isInternalLink(tt.args.url, tt.args.baseUrl); got != tt.want {
				t.Errorf("isInternalLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setHTMLVersion(t *testing.T) {
	type args struct {
		html string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"html5", args{html: sampleHTML5}, HTML5},
		{"html4_strict", args{html: sampleHTML4}, HTML401_STRICT},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setHTMLVersion(tt.args.html); got != tt.want {
				t.Errorf("setHTMLVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppData(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *AppData
	}{
		{"demo_website", args{url: "http://eshankavishka.tech/"}, &AppData{
			HtmlVersion:  HTML5,
			Title:        "Eshan Kavishka portfolio",
			HeadingCount: complexHtmlHeaderCount,
			Links:        complexHTMLLinks,
			HasLogin:     false,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAppData(tt.args.url)
			if got.HtmlVersion != tt.want.HtmlVersion {
				t.Errorf("GetAppData() : HTML version = %v, want %v", got.HtmlVersion, tt.want.HtmlVersion)
			}
			if !reflect.DeepEqual(got.HeadingCount, tt.want.HeadingCount) {
				t.Errorf("GetAppData() : Heading count = %v, want %v", got.HeadingCount, tt.want.HeadingCount)
			}
			if got.Title != tt.want.Title {
				t.Errorf("GetAppData() : Title = %v, want %v", got.Title, tt.want.Title)
			}
			if got.HasLogin != tt.want.HasLogin {
				t.Errorf("GetAppData() : Has login page = %v, want %v", got.HasLogin, tt.want.HasLogin)
			}
		})
	}
}
