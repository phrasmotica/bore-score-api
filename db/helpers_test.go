package db

import "testing"

func TestComputeName(t *testing.T) {
	tables := []struct {
		displayName string
		expected    string
	}{
		{"Uno", "uno"},
		{"Love Letter", "love-letter"},
		{"Odin's Ravens", "odins-ravens"},
		{"Uno: The Sequel", "uno-the-sequel"},
		{"Ganz Schön Clever", "ganz-schon-clever"},
		{"Sushi Go!", "sushi-go"},
	}

	for _, table := range tables {
		actual := computeName(table.displayName)
		if actual != table.expected {
			t.Errorf("Computed name was incorrect! Actual: %s, expected: %s", actual, table.expected)
		}
	}
}
