package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	// Add other fields as needed
}

type Mixed struct {
	Type string `json:"type"`
	Main []struct {
		Type  string `json:"type"`
		Index int    `json:"index"`
		All   bool   `json:"all"`
		// Add other fields as needed
	} `json:"main"`
	// Add other fields as needed
}

type Web struct {
	Type    string         `json:"type"`
	Results []SearchResult `json:"results"`
	// Add other fields as needed
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
	// Add other fields as needed
}

type MetaURL struct {
	Scheme   string `json:"scheme"`
	Netloc   string `json:"netloc"`
	Hostname string `json:"hostname"`
	Favicon  string `json:"favicon"`
	Path     string `json:"path"`
	// Add other fields as needed
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

	searchQuery := strings.Join(os.Args[1:], " ")
	url := fmt.Sprintf("https://api.search.brave.com/res/v1/web/search?q=site:%s", searchQuery)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Subscription-Token", braveAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON response into the Response struct
	var searchResponse Response
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	for _, result := range searchResponse.Web.Results {
		fmt.Println(result.URL)
	}
}
