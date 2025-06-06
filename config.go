package main

// this one is free API that does not require any API key
const ExternalAPIBaseURL = "https://jsonmock.hackerrank.com/api/football_matches?page="

func GetExternalAPIBaseURL() string {
	return ExternalAPIBaseURL
}
