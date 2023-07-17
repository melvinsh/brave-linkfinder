package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// Define the Go structs to match the JSON response
type Response struct {
	Query Query  `json:"query"`
	Mixed Mixed  `json:"mixed"`
	Type  string `json:"type"`
	Web   Web    `json:"web"`
}

type Query struct {
	Original          string `json:"original"`
	ShowStrictWarning bool   `json:"show_strict_warning"`
	IsNavigational    bool   `json:"is_navigational"`
}

type Mixed struct {
	Type string `json:"type"`
	Main []struct {
		Type  string `json:"type"`
		Index int    `json:"index"`
		All   bool   `json:"all"`
	} `json:"main"`
}

type Web struct {
	Type    string         `json:"type"`
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Title          string  `json:"title"`
	URL            string  `json:"url"`
	IsSourceLocal  bool    `json:"is_source_local"`
	IsSourceBoth   bool    `json:"is_source_both"`
	Description    string  `json:"description"`
	Language       string  `json:"language"`
	FamilyFriendly bool    `json:"family_friendly"`
	Type           string  `json:"type"`
	Subtype        string  `json:"subtype"`
	MetaURL        MetaURL `json:"meta_url"`
	Age            string  `json:"age"`
}

type MetaURL struct {
	Scheme   string `json:"scheme"`
	Netloc   string `json:"netloc"`
	Hostname string `json:"hostname"`
	Favicon  string `json:"favicon"`
	Path     string `json:"path"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a search query as a command-line argument.")
		return
	}

	braveAPIKey := os.Getenv("BRAVE_API_KEY")
	if braveAPIKey == "" {
		fmt.Println("Please set the BRAVE_API_KEY environment variable.")
		return
	}

	var results []SearchResult

	for i := 0; i < 9; i++ {
		// todo: error handling
		searchResponse, _ := request(braveAPIKey, 0)

		if len(searchResponse.Web.Results) == 0 {
			break
		}

		results = append(results, searchResponse.Web.Results...)

		// todo: make variable based on actual key limit
		time.Sleep(time.Second)
	}

	var urls []string

	for _, result := range results {
		urls = append(urls, result.URL)
	}

	urls = removeDuplicatesAndSort(urls)

	for _, url := range urls {
		fmt.Println(url)
	}
}

func request(braveAPIKey string, offset int) (Response, error) {
	searchQuery := strings.Join(os.Args[1:], " ")
	// offset of 9 and count of 20 are the current documented maximums
	// https://api.search.brave.com/app/documentation/query
	url := fmt.Sprintf("https://api.search.brave.com/res/v1/web/search?q=site:%s&count=%d&offset=%d", searchQuery, 20, offset)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return Response{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Subscription-Token", braveAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		return Response{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return Response{}, err
	}

	var searchResponse Response
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return Response{}, err
	}
	return searchResponse, nil
}

func removeDuplicatesAndSort(arr []string) []string {
	seen := make(map[string]bool)
	uniqueStrings := []string{}

	for _, str := range arr {
		if !seen[str] {
			seen[str] = true
			uniqueStrings = append(uniqueStrings, str)
		}
	}

	sort.Strings(uniqueStrings)

	return uniqueStrings
}
