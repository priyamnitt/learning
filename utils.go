package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// fetchAllMatches grabs match data from the external API for the first 20 pages
// just because the API has 1159 pages so loading whole data will take too much time.
func fetchAllMatches() ([]Match, error) {
	baseURL := GetExternalAPIBaseURL()
	var allMatches []Match

	for page := 1; page <= 20; page++ {
		fmt.Println("Fetching", baseURL+strconv.Itoa(page))
		resp, err := http.Get(baseURL + strconv.Itoa(page))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var apiResponse APIResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return nil, err
		}

		allMatches = append(allMatches, apiResponse.Data...)
	}

	return allMatches, nil
}

// calculateTeamStats builds a map of team names to their stats
// by iterating through all matches.
func calculateTeamStats(matches []Match) map[string]TeamStats {
	teamStats := make(map[string]TeamStats)

	for _, match := range matches {
		team1 := match.Team1
		team2 := match.Team2
		team1Goals, _ := strconv.Atoi(match.Team1Goals)
		team2Goals, _ := strconv.Atoi(match.Team2Goals)

		// Update stats for team1
		stats1 := teamStats[team1]
		stats1.Team = team1
		stats1.MatchesPlayed++
		stats1.GoalsFor += team1Goals
		stats1.GoalsAgainst += team2Goals
		stats1.GoalDifference = stats1.GoalsFor - stats1.GoalsAgainst

		if team1Goals > team2Goals {
			stats1.Wins++
			stats1.Points += 3
		} else if team1Goals == team2Goals {
			stats1.Draws++
			stats1.Points++
		} else {
			stats1.Losses++
		}
		teamStats[team1] = stats1

		// Update stats for team2
		stats2 := teamStats[team2]
		stats2.Team = team2
		stats2.MatchesPlayed++
		stats2.GoalsFor += team2Goals
		stats2.GoalsAgainst += team1Goals
		stats2.GoalDifference = stats2.GoalsFor - stats2.GoalsAgainst

		if team2Goals > team1Goals {
			stats2.Wins++
			stats2.Points += 3
		} else if team2Goals == team1Goals {
			stats2.Draws++
			stats2.Points++
		} else {
			stats2.Losses++
		}
		teamStats[team2] = stats2
	}

	return teamStats
}
