package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var googleDomains = map[string]string{
	"com": "https://www.google.com/search?q=",
}

type SearchResult struct {
	ResultRank  int
	ResultUrl   string
	ResultTitle string
	ResultDesc  string
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func buildGoogleUrls(searchTerm string, languageCode string, countryCode string, pages int, count int) ([]string, error) {
	toScrape := []string{}
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBase, found := googleDomains[countryCode]; found {
		for i := 0; i < pages; i++ {
			start := i * count
			scrapeURL := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0", googleBase, searchTerm, count, languageCode, start)
			toScrape = append(toScrape, scrapeURL)

		}
	} else {
		err := fmt.Errorf("country (%s) is currently not supported", countryCode)
		return nil, err
	}
	return toScrape, nil
}

func GoogleScrape(searchTem string, languageCode string, countryCode string, pages int, count int) ([]SearchResult, error) {
	results := []SearchResult{}
	resultCounter := 0
	googlePages, err := buildGoogleUrls(searchTem, languageCode, countryCode, pages, count)

	if err != nil {
		return nil, err
	}

	for _, page := range googlePages {
		res, err := ScrapeClientRequest(page)

		if err != nil {
			return nil, err
		}

		data, err := googleResultParsing(res, resultCounter)

		if err != nil {
			return nil, err
		}

		resultCounter += len(data)

		for _, result := range data {
			results = append(results, result)
		}
		time.Sleep(time.Second * 10)
	}

	return results, nil
}

func ScrapeClientRequest(searchURL string) (*http.Response, error) {
	baseClient := getScrapeClient()
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := baseClient.Do(req)

	if res.StatusCode != 200 {
		err := fmt.Errorf("Scraper received a non-200 code")
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func getScrapeClient() *http.Client {
	return &http.Client{}
}

func googleResultParsing() {

}

func main() {
	res, err := GoogleScrape("michalis anastasiou", "en", "com", 1, 30)
	if err == nil {
		for _, res := range res {
			fmt.Println(res)
		}
	}
}
