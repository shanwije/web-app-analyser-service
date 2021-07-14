package collector

const (
	UNKNOWN              = "UNKNOWN"
	HTML5                = "HTML 5"
	XHTML11              = "XHTML 1.1"
	XHTML11_FRAMESET     = "XHTML 1.0 Frameset"
	XHTML11_TRANSITIONAL = "XHTML 1.0 Transitional"
	XHTML10_STRICT       = "XHTML 1.0 Strict"
	HTML401_FRAMESET     = "HTML 4.01 Frameset"
	HTML401_TRANSITIONAL = "HTML 4.01 Transitional"
	HTML401_STRICT       = "HTML 4.01 Strict"
)

const (
	H1 = "h1"
	H2 = "h2"
	H3 = "h3"
	H4 = "h4"
	H5 = "h5"
	H6 = "h6"
)

func GetHtmlVersions() map[string]string {
	return map[string]string{
		HTML401_STRICT:       `"-//W3C//DTD HTML 4.01//EN"`,
		HTML401_TRANSITIONAL: `"-//W3C//DTD HTML 4.01 Transitional//EN"`,
		HTML401_FRAMESET:     `"-//W3C//DTD HTML 4.01 Frameset//EN"`,
		XHTML10_STRICT:       `"-//W3C//DTD XHTML 1.0 Strict//EN"`,
		XHTML11_TRANSITIONAL: `"-//W3C//DTD XHTML 1.0 Transitional//EN"`,
		XHTML11_FRAMESET:     `"-//W3C//DTD XHTML 1.0 Frameset//EN"`,
		XHTML11:              `"-//W3C//DTD XHTML 1.1//EN"`,
		HTML5:                `<!DOCTYPE html>`,
	}
}

const COLLY_TIMEOUT_DURATION = 10
const LINK_LIST_COLLECTOR_THREAD_COUNT = 4
const LINK_LIST_COLLECTOR_DEPTH = 2
