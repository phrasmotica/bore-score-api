package routes

import "phrasmotica/bore-score-api/models"

func hasUniquePlayerScores(result *models.Result) bool {
	var uniquePlayers []string

	for _, e := range result.Scores {
		// TODO: use slices.Contains(...) to check
		uniquePlayers = appendIfMissing(uniquePlayers, e.Username)
	}

	return len(uniquePlayers) == len(result.Scores)
}

func appendIfMissing(slice []string, s string) []string {
	for _, e := range slice {
		if e == s {
			return slice
		}
	}

	return append(slice, s)
}
