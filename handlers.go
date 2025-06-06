package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Returns the top 3 teams by points across all years
func getTopTeams(res http.ResponseWriter, req *http.Request) {

	matches, err := fetchAllMatches()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	teamStats := calculateTeamStats(matches)

	var teams []TeamStats
	for _, stats := range teamStats {
		teams = append(teams, stats)
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Points > teams[j].Points
	})

	topN := 3
	if len(teams) < 3 {
		topN = len(teams)
	}
	topTeams := teams[:topN]

	json.NewEncoder(res).Encode(map[string]interface{}{
		"success": true,
		"data":    topTeams,
	})
}

// Returns the team with the most goals scored till now
func getMaxGoals(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	matches, err := fetchAllMatches()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	teamGoals := make(map[string]int)

	for _, match := range matches {
		team1Goals, _ := strconv.Atoi(match.Team1Goals)
		team2Goals, _ := strconv.Atoi(match.Team2Goals)

		teamGoals[match.Team1] += team1Goals
		teamGoals[match.Team2] += team2Goals
	}

	var maxTeam string
	var maxGoals int
	for team, goals := range teamGoals {
		if goals > maxGoals {
			maxGoals = goals
			maxTeam = team
		}
	}

	json.NewEncoder(res).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"team":  maxTeam,
			"goals": maxGoals,
		},
	})
}

// Returns the total number of drawn matches across all years
func getDraws(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	matches, err := fetchAllMatches()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	draws := 0
	for _, match := range matches {
		if match.Team1Goals == match.Team2Goals {
			draws++
		}
	}
	json.NewEncoder(res).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"draws": draws,
		},
	})
}

// Returns the total number of goals scored in all matches across all teams and all years till now
func getTotalGoals(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	matches, err := fetchAllMatches()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	totalGoals := 0
	for _, match := range matches {
		team1Goals, _ := strconv.Atoi(match.Team1Goals)
		team2Goals, _ := strconv.Atoi(match.Team2Goals)
		totalGoals += team1Goals + team2Goals
	}
	json.NewEncoder(res).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"total_goals": totalGoals,
		},
	})
}

// Handles /api/teams/ routes, dispatches to team stats or head-to-head based on the path
func teamsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	path := req.URL.Path[len("/api/teams/"):]
	if strings.Contains(path, "/vs/") {
		getHeadToHead(res, req, path)
	} else {
		getTeamStats(res, req, path)
	}
}

// Returns stats for a single team as provided by the user as per the path
func getTeamStats(res http.ResponseWriter, req *http.Request, teamName string) {
	matches, err := fetchAllMatches()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	teamStats := calculateTeamStats(matches)
	stats, exists := teamStats[teamName]
	if !exists {
		http.Error(res, "Team not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(res).Encode(map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// Returns head-to-head stats between two teams as per the request
func getHeadToHead(res http.ResponseWriter, req *http.Request, path string) {
	parts := strings.Split(path, "/vs/")
	if len(parts) != 2 {
		http.Error(res, "Invalid path format. Use /api/teams/team1/vs/team2", http.StatusBadRequest)
		return
	}
	team1 := parts[0]
	team2 := parts[1]
	matches, err := fetchAllMatches()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	var headToHeadMatches []Match
	team1Wins := 0
	team2Wins := 0
	draws := 0
	team1Goals := 0
	team2Goals := 0
	for _, match := range matches {
		if (match.Team1 == team1 && match.Team2 == team2) || (match.Team1 == team2 && match.Team2 == team1) {
			headToHeadMatches = append(headToHeadMatches, match)
			t1g, _ := strconv.Atoi(match.Team1Goals)
			t2g, _ := strconv.Atoi(match.Team2Goals)
			if match.Team1 == team1 {
				team1Goals += t1g
				team2Goals += t2g
				if t1g > t2g {
					team1Wins++
				} else if t1g < t2g {
					team2Wins++
				} else {
					draws++
				}
			} else {
				team1Goals += t2g
				team2Goals += t1g
				if t2g > t1g {
					team1Wins++
				} else if t2g < t1g {
					team2Wins++
				} else {
					draws++
				}
			}
		}
	}
	json.NewEncoder(res).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"team1":       team1,
			"team2":       team2,
			"matches":     len(headToHeadMatches),
			"team1_wins":  team1Wins,
			"team2_wins":  team2Wins,
			"draws":       draws,
			"team1_goals": team1Goals,
			"team2_goals": team2Goals,
		},
	})
}
