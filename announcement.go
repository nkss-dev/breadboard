package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "regexp"

    "github.com/PuerkitoBio/goquery"
)

type Announcement struct {
    Name string   `json:"name"`
    Link string   `json:"link"`
    Tags []string `json:"tags"`
}

func check(err error) {
    if err != nil {
        fmt.Println(err)
    }
}

func FormatAnnouncement(name string) []string {
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
    url := "http://nitkkr.ac.in/notifications.php?tbl=notifications"

    response, err := http.Get(url)
    check(err)
    defer response.Body.Close()

    doc, err := goquery.NewDocumentFromReader(response.Body)
    check(err)

    var announcements []Announcement
    doc.Find("div.bg-white").Find("p").Each(func(index int, item *goquery.Selection) {
        var announcement Announcement
        announcement.Name = item.Find("a").Text()
        announcement.Link, _ = item.Find("a").Attr("href")
        announcement.Tags = FormatAnnouncement(announcement.Name)

        if announcement.Name != "" && announcement.Link != "" {
            announcements = append(announcements, announcement)
        }
    })

    jsonResp, _ := json.Marshal(announcements)
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResp)
}
