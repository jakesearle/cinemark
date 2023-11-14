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
const baseUrl = "https://www.cinemark.com/theatres"
const defaultRegion = "nc-raleigh"
const defaultTheatre = "cinemark-bistro-raleigh"

func ScrapeFunc() {
	thisWeeksMovies()
}

func thisWeeksMovies() {

	currentDate := time.Now().Format("2006-01-02")
	fmt.Println(buildUrl(defaultRegion, defaultTheatre, currentDate))

	// requestUrl := baseUrl + authorPath
	// doc := GetSoup(requestUrl)
	// authorNodes := QueryAll(doc, ".body ul li")
	// authors := make([]*Author, 0)
	// for _, authorNode := range authorNodes {
	// 	pTags := QueryAll(authorNode, "p")
	// 	if len(pTags) < 1 {
	// 		fmt.Println("There's no author name here... Hm.")
	// 		continue
	// 	}
	// 	if len(pTags) < 2 {
	// 		fmt.Println("This author has no citations")
	// 	}
	// 	authorName := TrimLeadingAsterisks(GetText(pTags[0]))
	// 	credits := Map2(pTags[1:], GetInt)
	// 	author := &Author{
	// 		Name:    authorName,
	// 		Credits: credits,
	// 	}
	// 	authors = append(authors, author)
	// }
	// return authors
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
	pathedUrl := fmt.Sprintf("%s/%s/%s", baseUrl, regionName, theatreName)
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
