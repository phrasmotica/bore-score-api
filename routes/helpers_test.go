package routes

import (
	"phrasmotica/bore-score-api/models"
	"testing"
)

func TestHasUniquePlayerScores(t *testing.T) {
	tables := []struct {
		result   models.Result
		expected bool
	}{
		{models.Result{Scores: make([]models.PlayerScore, 0)}, true},

		{models.Result{
			Scores: []models.PlayerScore{
				{
					Username: "player1",
				},
			},
		}, true},

		{models.Result{
			Scores: []models.PlayerScore{
				{
					Username: "player1",
				},
				{
					Username: "player2",
				},
			},
		}, true},

		{models.Result{
			Scores: []models.PlayerScore{
				{
					Username: "player2",
				},
				{
					Username: "player2",
				},
			},
		}, false},

		{models.Result{
			Scores: []models.PlayerScore{
				{
					Username: "player1",
				},
				{
					Username: "player2",
				},
				{
					Username: "player2",
				},
			},
		}, false},
	}

	for _, table := range tables {
		actual := hasUniquePlayerScores(&table.result)
		if actual != table.expected {
			t.Errorf("Computed value was incorrect! Actual: %t, expected: %t", actual, table.expected)
		}
	}
}
