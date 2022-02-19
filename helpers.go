package main

import "phrasmotica/bore-score-api/models"

// returns the highest ID in the given list of results
func getMaxResultId(results []models.Result) int {
	var maxId int

	for i, e := range results {
		if i == 0 || e.ID > maxId {
			maxId = e.ID
		}
	}

	return maxId
}

// returns the highest ID in the given list of players
func getMaxPlayerId(players []models.Player) int {
	var maxId int

	for i, e := range players {
		if i == 0 || e.ID > maxId {
			maxId = e.ID
		}
	}

	return maxId
}

// returns the highest ID in the given list of games
func getMaxGameId(games []models.Game) int {
	var maxId int

	for i, e := range games {
		if i == 0 || e.ID > maxId {
			maxId = e.ID
		}
	}

	return maxId
}
