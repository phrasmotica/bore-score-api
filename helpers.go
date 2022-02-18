package main

// removes the player with the given username from the given list of players
func removePlayer(players []player, username string) []player {
	var indexOf int

	for i, p := range players {
		if username == p.Username {
			indexOf = i
			break
		}
	}

	return append(players[:indexOf], players[indexOf+1:]...)
}

// removes the game with the given ID from the given list of games
func removeGame(games []game, id int) []game {
	var indexOf int

	for i, g := range games {
		if id == g.ID {
			indexOf = i
			break
		}
	}

	return append(games[:indexOf], games[indexOf+1:]...)
}

// removes all results with the given game ID from the given list of results
func removeResultsOfGame(results []result, gameId int) []result {
	remaining := []result{}

	for i := range results {
		if results[i].GameID != gameId {
			remaining = append(remaining, results[i])
		}
	}

	return remaining
}

// returns whether the game with the given ID exists in the given list of games
func gameExists(games []game, gameId int) bool {
	for _, g := range games {
		if g.ID == gameId {
			return true
		}
	}

	return false
}

// returns whether the player with the given ID exists in the given list of players
func playerExists(players []player, playerId int) bool {
	for _, p := range players {
		if p.ID == playerId {
			return true
		}
	}

	return false
}

// returns whether the player with the given username exists in the given list of players
func playerExistsByUsername(players []player, username string) bool {
	for _, p := range players {
		if p.Username == username {
			return true
		}
	}

	return false
}

// returns the highest ID in the given list of results
func getMaxResultId(results []result) int {
	var maxId int

	for i, e := range results {
		if i == 0 || e.ID > maxId {
			maxId = e.ID
		}
	}

	return maxId
}

// returns the highest ID in the given list of players
func getMaxPlayerId(players []player) int {
	var maxId int

	for i, e := range players {
		if i == 0 || e.ID > maxId {
			maxId = e.ID
		}
	}

	return maxId
}

// returns the highest ID in the given list of games
func getMaxGameId(games []game) int {
	var maxId int

	for i, e := range games {
		if i == 0 || e.ID > maxId {
			maxId = e.ID
		}
	}

	return maxId
}
