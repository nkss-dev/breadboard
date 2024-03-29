package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	query "breadboard/.sqlc-auto-gen"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/html"
)

type Announcement struct {
	Date  string   `json:"date"`
	Title string   `json:"title"`
	Link  string   `json:"link"`
	Tags  []string `json:"tags"`
}

const query_prefix = "INSERT INTO academic_announcement (date_of_creation, title, title_link, kind) VALUES "

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

// scrapeAnnouncements returns all the announcements from a specified URL.
//
// It scrapes the URL and retrieves elements
// to convert them into the Announcement type.
func scrapeAnnouncements() (announcements []Announcement) {
	// Request the HTML page
	response, err := http.Get("https://nitkkr.ac.in/?page_id=621")
	if err != nil {
		//RespondError(w, 404, "The source web-page for scraping was not found")
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		//RespondError(w, response.StatusCode, "")
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		//RespondError(w, 502, "Server could not parse the source HTML document")
		return
	}

	// Find the announcements
	doc.Find("div.comman-inner-section").Find("p").Each(func(i int, item *goquery.Selection) {
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

// parseDate() attempts to convert obtained date to time.Time
//
// It also has some case-specific checks since scraped data doesn't have
// uniformity
func parseDate(date string) (parsedDate pgtype.Date, err error) {
	date = strings.Replace(date, ".", "-", 3)
	date = strings.Replace(date, "/", "-", 3)
	date = strings.Replace(date, "__", "", 2)
	if date == "1-12-2021" {
		date = "01-12-2021"
	}
	if date == "07-05-019" {
		date = "07-05-2019"
	}
	time, err := time.Parse("02-01-2006", date)
	return pgtype.Date{Time: time}, err
}

// insertNewAnnouncements saves announcements to database
//
// It first calls scrapeAnnouncements() and then saves the results after
// some formatting
//
// TODO: try to use sqlc or some other intermediate to store this query
func FetchAnnouncements(pool *pgxpool.Pool) {
	insert_query := query_prefix
	ctx := context.Background()
	queries := query.New(pool)
	latest_date, err := queries.GetLatestAnnouncementDate(ctx)
	if err != nil {
		fmt.Println(err)
	}

	var announcements = scrapeAnnouncements()
	for _, announcement := range announcements {
		date, err := parseDate(announcement.Date)
		if err != nil {
			fmt.Println(date, err)
		} else {
			if err != nil && !date.Time.After(latest_date.Time) {
				fmt.Println("Detected old announcement, assuming it is in database already")
				break
			}
			addition := strings.Join([]string{
				"('",
				date.Time.Format("2006-01-02"),
				"', '",
				announcement.Title,
				"', '",
				announcement.Link,
				"',  'academic') ",
			}, "")
			if insert_query != query_prefix {
				addition = ", " + addition
			}
			insert_query += addition
		}
	}
	insert_query += " ON CONFLICT (date_of_creation, title) DO NOTHING"
	_, inserterr := pool.Exec(ctx, insert_query)
	if inserterr != nil {
		fmt.Println(inserterr)
	}
}

// GetAnnouncements returns all the announcements stored in database
//
// It is a wrapper function around
func GetAnnouncements(pool *pgxpool.Pool) http.HandlerFunc {
	ctx := context.Background()
	queries := query.New(pool)
	return func(w http.ResponseWriter, r *http.Request) {
		announcements, err := queries.GetAcademicAnnouncements(ctx)
		if err == pgx.ErrNoRows || len(announcements) == 0 {
			FetchAnnouncements(pool)
		}
		announcements, err = queries.GetAcademicAnnouncements(ctx)
		if err == pgx.ErrNoRows {
			RespondError(w, 404, "Announcements not found in the database")
			return
		}
		RespondJSON(w, 200, announcements)
	}
}
