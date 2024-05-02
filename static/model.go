package static

// Season contains an ID and a name for a season
type Season struct {
	ID     int64  `json:"id"`
	Season string `json:"season"`
}

// Queue contains a description and notes, and ID and a map for a queue
type Queue struct {
	ID          int64  `json:"queueId"`
	Map         string `json:"map"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
}

// Map contains notes, an ID and a name for a map
type Map struct {
	ID    int64  `json:"mapId"`
	Name  string `json:"mapName"`
	Notes string `json:"notes"`
}

// GameMode contains a description and name of a game mode
type GameMode struct {
	Mode        string `json:"gameMode"`
	Description string `json:"description"`
}

// GameType contains a description and name of a game type
type GameType struct {
	Type        string `json:"gameType"`
	Description string `json:"description"`
}
