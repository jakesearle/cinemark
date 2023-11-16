package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// const completeUrl = "https://www.cinemark.com/theatres/nc-raleigh/cinemark-bistro-raleigh?showDate=2023-11-15"
const baseUrl = "https://www.cinemark.com"
const theatreUrl = baseUrl + "/theatres"
const defaultRegion = "nc-raleigh"
const defaultTheatre = "cinemark-bistro-raleigh"

func ScrapeFunc() {
	calendarWeekFilms()
}

func calendarWeekFilms() {
	currentTime := time.Now()
	lastDay := currentTime.AddDate(0, 0, 6-int(currentTime.Weekday()))

	films := make(map[Film]bool)
	for iterDay := currentTime; iterDay.Before(lastDay) || isSameDay(iterDay, lastDay); iterDay = iterDay.AddDate(0, 0, 1) {
		dateStr := fmt.Sprintf("%s", iterDay.Format("2006-01-02"))
		appendUniqueFilms(dateStr, &films)
	}
	for k := range films {
		fmt.Println(k)
	}
}

func appendUniqueFilms(currentDate string, allFilms *map[Film]bool) {
	// Date in the format "2006-01-02"
	fmt.Printf("%s: %d", currentDate, len(*allFilms))
	url := buildUrl(defaultRegion, defaultTheatre, currentDate)
	soup := GetSoup(url)
	filmNodes := QueryAll(soup, "#showTimes .showtimeMovieBlock")
	for _, filmNode := range filmNodes {
		film := Film{
			Title:     getTitle(filmNode),
			Link:      getLink(filmNode),
			PosterUrl: getPosterUrl(filmNode),
		}
		(*allFilms)[film] = true
	}
	fmt.Printf(" -> %d\n", len(*allFilms))
}

func getTitle(filmNode *html.Node) string {
	return GetText(Query(filmNode, ".movieLink h3"))
}

func getLink(filmNode *html.Node) string {
	linkNode := Query(filmNode, ".movieLink")
	path := AttrOr(linkNode, "href", "")
	if path == "" {
		return ""
	}
	return baseUrl + path
}

func getPosterUrl(filmNode *html.Node) string {
	pictureNode := Query(filmNode, "picture img")
	return AttrOr(pictureNode, "data-srcset", "")
}

func GetSoup(url string) *html.Node {
	htmlString := LoadOrCacheHtml(url)
	soup, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		log.Fatal(err)
	}
	return soup
}

func buildUrl(regionName, theatreName, date string) string {
	// Creating a URL object
	pathedUrl := fmt.Sprintf("%s/%s/%s", theatreUrl, regionName, theatreName)
	u, err := url.Parse(pathedUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		panic(err)
	}
	// Adding query parameters
	q := u.Query()
	q.Set("showDate", date)
	u.RawQuery = q.Encode()
	return u.String()
}

func isSameDay(time1, time2 time.Time) bool {
	return time1.Year() == time2.Year() &&
		time1.Month() == time2.Month() &&
		time1.Day() == time2.Day()
}
