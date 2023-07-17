# Brave Search API Linkfinder

This is a simple Go program that interacts with the Brave Search API to find links related to a certain domain.

## Prerequisites

Before running this program, you need to set up the following:

1. Go programming language installed on your system.
2. A valid Brave Search API key. If you don't have one, you can obtain one [here](https://api.search.brave.com/app/keys).

## Usage

To run the Go program, follow these steps:

Set the `BRAVE_API_KEY` environment variable with your valid Brave Search API key:

```bash
export BRAVE_API_KEY=your_api_key_here
```

Execute the program with your desired search query as the command-line argument:
``` bash
go run main.go hadrian.io
```

The program will send a request to the Brave Search API and display the parsed response data. The output will be all the URLs from the web search results retrieved from the API. For example:

```
https://hadrian.io/company/about-hadrian
https://hadrian.io/company/careers
https://hadrian.io
https://hadrian.io/blog/non-techies-in-a-tech-world-how-hadrian-makes-learning-about-tech-accessible-not-daunting
https://hadrian.io/blog/technology-is-changing-and-the-security-strategy-needs-to-change-with-it
https://hadrian.io/blog/mimicking-a-hacker-with-event-based-ai
```