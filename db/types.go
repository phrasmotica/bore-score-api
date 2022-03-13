package db

type Summary struct {
	GameCount   int64 `json:"gameCount"`
	GroupCount  int64 `json:"groupCount"`
	PlayerCount int64 `json:"playerCount"`
	ResultCount int64 `json:"resultCount"`
}
