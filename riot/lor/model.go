package lor

// Player represents a ranked Legends of Runeterra player.
type Player struct {
	Name         string `json:"name"`
	Rank         int64  `json:"rank"`
	LeaguePoints int64  `json:"lp"`
}
