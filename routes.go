package main

import "net/http"

func registerRoutes() {
	http.HandleFunc("/api/tournament/top-teams", getTopTeams)
	http.HandleFunc("/api/tournament/max-goals", getMaxGoals)
	http.HandleFunc("/api/tournament/draws", getDraws)
	http.HandleFunc("/api/tournament/total-goals", getTotalGoals)
	http.HandleFunc("/api/teams/", teamsHandler)
}
