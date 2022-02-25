package db

type Summary struct {
	GameCount   int64 `json:"gameCount"`
	PlayerCount int64 `json:"playerCount"`
	ResultCount int64 `json:"resultCount"`
}
