package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type Announcement struct {
	Title string   `json:"title"`
	Link  string   `json:"link"`
	Tags  []string `json:"tags"`
}

func fetchTags(name string) []string {
	var tags []string
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

func GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	// Request the HTML page
	response, err := http.Get("http://nitkkr.ac.in/notifications.php")
	if err != nil {
		respondError(w, 404, "The source web-page for scraping was not found")
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		respondError(w, response.StatusCode, "")
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		respondError(w, 502, "Server could not parse the source HTML document")
		return
	}

	// Find the announcements
	var announcements []Announcement
	doc.Find("div.bg-white").Find("p").Each(func(i int, item *goquery.Selection) {
		title := item.Find("a").Text()
		link, _ := item.Find("a").Attr("href")
		encoded_url, err := url.Parse(link)
		if err != nil {
			fmt.Println(err)
			return
		}
		link = encoded_url.String()
		tags := fetchTags(title)

		if title != "" && link != "" {
			announcements = append(announcements, Announcement{Title: title, Link: link, Tags: tags})
		}
	})

	respondJSON(w, 200, announcements)
}
