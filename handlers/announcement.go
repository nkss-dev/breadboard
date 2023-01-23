package handlers

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Announcement struct {
	Date  string   `json:"date"`
	Title string   `json:"title"`
	Link  string   `json:"link"`
	Tags  []string `json:"tags"`
}

// fetchTags returns the tags for a given string.
//
// It retrieves tags from the given string using RegEx for each
// tag that is to be searched for. Due to its high complexity
// and unnecessary computation, it will be replaces in the future.
func fetchTags(name string) (tags []string) {
	regexes := map[string]*regexp.Regexp{
		// Disciplines
		"CE":  regexp.MustCompile(`(?i)( CE )|(civil engg[\. ]? ?deptt)`),
		"CS":  regexp.MustCompile(`(?i)( CS )|(Computer( Engineering)?( Department)?)`),
		"ECE": regexp.MustCompile(`(?i)( ECE )|(Electronics (and|&) Communication( Department)?)`),
		"EE":  regexp.MustCompile(`(?i)( EE )|(Electrical( Engineering)?( Department)?)`),
		"ME":  regexp.MustCompile(`(?i)( ME )|(Mechanical( Engineering)?( Department)?)`),
		"IT":  regexp.MustCompile(`(?i)( IT )|(Information Technology( Department)?)`),
		"PIE": regexp.MustCompile(`(?i)( PIE )|(Production ((and|&) )(Industrial )?Engineering( Department)?)`),

		// Degree
		"B.Tech.": regexp.MustCompile(`(?i)B[\. ]?Tech`),
		"M.Tech.": regexp.MustCompile(`(?i)M[\. ]?Tech`),
		"MCA":     regexp.MustCompile(`(?i)MCA`),
		"MBA":     regexp.MustCompile(`(?i)MBA`),
		"Ph.D":    regexp.MustCompile(`(?i)Ph[\. ]?D\.?`),

		// Semester
		"1st semester": regexp.MustCompile(`(?i)[^(except )]1st sem(ester)?`),
		"2nd semester": regexp.MustCompile(`(?i)[^(except )]2nd sem(ester)?`),
		"3rd semester": regexp.MustCompile(`(?i)[^(except )]3rd sem(ester)?`),
		"4th semester": regexp.MustCompile(`(?i)[^(except )]4th sem(ester)?`),
		"5th semester": regexp.MustCompile(`(?i)[^(except )]5th sem(ester)?`),
		"6th semester": regexp.MustCompile(`(?i)[^(except )]6th sem(ester)?`),
		"7th semester": regexp.MustCompile(`(?i)[^(except )]7th sem(ester)?`),
		"8th semester": regexp.MustCompile(`(?i)[^(except )]8th sem(ester)?`),

		// Examination
		"Mid Sem":     regexp.MustCompile(`(?i)mid sem(ester)?`),
		"Mid Sem - 1": regexp.MustCompile(`(?i)mid sem(ester)? (exam|test)-I[^I]`),
		"Mid Sem - 2": regexp.MustCompile(`(?i)mid sem(ester)? (exam|test)-II`),
		"End Sem":     regexp.MustCompile(`(?i)end sem(ester)? (exam|test)`),
	}

	for key, value := range regexes {
		if tag := value.Find([]byte(name)); tag != nil {
			tags = append(tags, key)
		}
	}
	return tags
}

// getTextInSpan returns the textual data within a span tag.
//
// It iteratively traverses through each child within the span tag
// to search for `html.TextNode` and returns the trimmed data within.
func getTextInSpan(a *html.Node) (text string) {
	// Loop into the child continuously until a text node is found
	for n := a.FirstChild; n != nil; n = n.FirstChild {
		if n.Type == html.TextNode {
			text = strings.TrimSpace(n.Data)
			break
		}
	}
	return text
}

// parseA returns the text and its URL from a hyperlink.
//
// It looks for the href attribute within the `<a>` tag's attributes
// to find its text and URL. If it encounters `<span>`, it then calls
// `getTextInSpan` to retrieve the text within.
func parseA(a *html.Node) (text string, link string) {
	for _, attr := range a.Attr {
		if attr.Key != "href" {
			continue
		}

		encoded_url, _ := url.Parse(attr.Val)
		link = encoded_url.String()
	}

	if a.FirstChild.Data == "span" {
		text = getTextInSpan(a.FirstChild)
	} else {
		text = strings.TrimSpace(a.FirstChild.Data)
	}
	return text, link
}

// parseSpan returns the text and its URL from a hyperlink.
//
// It iteratively traveses through each child within the span tag
// to search for the `<a>` tag and calls `parseA` to retrieve the
// text and URL of the respective tag.
func parseSpan(a *html.Node) (text string, link string) {
	for n := a.FirstChild; n != nil; n = n.FirstChild {
		if n.Data == "a" {
			text, link = parseA(n)
			break
		}
	}
	return text, link
}

//scrapes announcements from official NIT website
func scrapeAnnouncements() (announcements []Announcement) {
	// Request the HTML page
    response, err := http.Get("https://nitkkr.ac.in/?page_id=621")
	if err != nil {
		//RespondError(w, 404, "The source web-page for scraping was not found")
		return announcements
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		//RespondError(w, response.StatusCode, "")
		return announcements
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		//RespondError(w, 502, "Server could not parse the source HTML document")
		return announcements
	}

	// Find the announcements
	doc.Find("div.bg-white").Find("p").Each(func(i int, item *goquery.Selection) {
		for _, node := range item.Nodes {
			for n := node.FirstChild; n != nil; n = n.NextSibling {
				if n.Type != html.ElementNode {
					continue
				}
				if n.Data != "a" && n.Data != "span" {
					continue
				}

				// Loop to previous siblings until a text node is found
				var date, PrevSibling string
				for prev := n.PrevSibling; prev != nil; prev = prev.PrevSibling {
					if prev.Data == "span" {
						PrevSibling = getTextInSpan(prev)
					} else if prev.Type == html.TextNode {
						PrevSibling = prev.Data
					}

					if PrevSibling != "" {
						break
					}
				}
				date = strings.TrimSpace(PrevSibling)

				var title, link string
				if n.Data == "a" {
					title, link = parseA(n)
				} else if n.Data == "span" {
					title, link = parseSpan(n)
				}

				tags := fetchTags(title)
				if title != "" && link != "" {
					announcements = append(announcements, Announcement{Date: date, Title: title, Link: link, Tags: tags})
				}
			}
		}
	})

	return announcements
}

// GetAnnouncements returns all the announcements from a specified URL.
//
// It is a handler function which scrapes the URL and retrieves elements
// to convert them into the Announcement type.
func GetAnnouncements() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	    var announcements = scrapeAnnouncements()
        RespondJSON(w, 200, announcements)
	}
}
