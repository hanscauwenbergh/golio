package lol

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/KnutZuidema/golio/datadragon"
	"github.com/KnutZuidema/golio/static"
)

// ChampionInfo contains information about the free champion rotation
type ChampionInfo struct {
	FreeChampionIDsForNewPlayers []int64 `json:"freeChampionIDsForNewPlayers"`
	FreeChampionIDs              []int64 `json:"freeChampionIDs"`
	MaxNewPlayerLevel            int64   `json:"maxNewPlayerLevel"`
}

// GetChampionsForNewPlayers returns data for champions available for free to new players
func (i *ChampionInfo) GetChampionsForNewPlayers(client *datadragon.Client) ([]datadragon.ChampionDataExtended, error) {
	res := make([]datadragon.ChampionDataExtended, 0, len(i.FreeChampionIDsForNewPlayers))
	for _, id := range i.FreeChampionIDsForNewPlayers {
		champion, err := client.GetChampionByID(strconv.FormatInt(id, 10))
		if err != nil {
			return nil, err
		}
		res = append(res, champion)
	}
	return res, nil
}

// GetChampions returns data for champions available for free
func (i *ChampionInfo) GetChampions(client *datadragon.Client) ([]datadragon.ChampionDataExtended, error) {
	res := make([]datadragon.ChampionDataExtended, 0, len(i.FreeChampionIDsForNewPlayers))
	for _, id := range i.FreeChampionIDs {
		champion, err := client.GetChampionByID(strconv.FormatInt(id, 10))
		if err != nil {
			return nil, err
		}
		res = append(res, champion)
	}
	return res, nil
}

// ChampionMastery represents the mastery of a champion in the mastery system for a summoner
type ChampionMastery struct {
	ChestGranted                 bool   `json:"chestGranted"`
	ChampionLevel                int64  `json:"championLevel"`
	ChampionPoints               int64  `json:"championPoints"`
	ChampionID                   int64  `json:"championId"`
	ChampionPointsUntilNextLevel int64  `json:"championPointsUntilNextLevel"`
	LastPlayTime                 int64  `json:"lastPlayTime"`
	TokensEarned                 int64  `json:"tokensEarned"`
	ChampionPointsSinceLastLevel int64  `json:"championPointsSinceLastLevel"`
	SummonerID                   string `json:"summonerId"`
}

// GetSummoner returns the summoner of this mastery
func (m *ChampionMastery) GetSummoner(client *Client) (*Summoner, error) {
	return client.Summoner.GetByID(m.SummonerID)
}

// GetChampion returns the champion of this mastery
func (m *ChampionMastery) GetChampion(client *datadragon.Client) (datadragon.ChampionDataExtended, error) {
	return client.GetChampionByID(strconv.FormatInt(m.ChampionID, 10))
}

// LeagueList represents a league containing all player entries in it
type LeagueList struct {
	LeagueID      string        `json:"leagueId"`
	Tier          string        `json:"tier"`
	Entries       []*LeagueItem `json:"entries"`
	Queue         string        `json:"queue"`
	Name          string        `json:"name"`
	sortedEntries []*LeagueItem
}

// GetRank returns the entry at the given rank, sorted by league points
func (l *LeagueList) GetRank(i int64) *LeagueItem {
	if l.sortedEntries == nil || len(l.sortedEntries) != len(l.Entries) {
		l.sortedEntries = make([]*LeagueItem, len(l.Entries))
		copy(l.sortedEntries, l.Entries)
		sort.Slice(
			l.sortedEntries, func(i, j int) bool {
				return l.sortedEntries[i].LeaguePoints > l.sortedEntries[j].LeaguePoints
			},
		)
	}
	return l.sortedEntries[i]
}

// LeagueItem represents a summoners ranked position in a league
type LeagueItem struct {
	QueueType    string      `json:"queueType"`
	SummonerName string      `json:"summonerName"`
	HotStreak    bool        `json:"hotStreak"`
	MiniSeries   *MiniSeries `json:"miniSeries"`
	Wins         int64       `json:"wins"`
	Veteran      bool        `json:"veteran"`
	Losses       int64       `json:"losses"`
	FreshBlood   bool        `json:"freshBlood"`
	Inactive     bool        `json:"inactive"`
	Tier         string      `json:"tier"`
	Rank         string      `json:"rank"`
	SummonerID   string      `json:"summonerId"`
	LeaguePoints int64       `json:"leaguePoints"`
}

// GetSummoner returns the summoner of this league item
func (i *LeagueItem) GetSummoner(client *Client) (*Summoner, error) {
	return client.Summoner.GetByID(i.SummonerID)
}

// MiniSeries represents a mini series when playing to ascend to the next ranked tier
type MiniSeries struct {
	Progress string `json:"progress"`
	Losses   int64  `json:"losses"`
	Target   int64  `json:"target"`
	Wins     int64  `json:"wins"`
}

// Match contains information about a match
type Match struct {
	// Match metadata
	Metadata *MatchMetadata `json:"metadata"`
	// Match info
	Info *MatchInfo `json:"info"`
}

// MatchMetadata contains metadata for a specific match
type MatchMetadata struct {
	// Match data version
	DataVersion string `json:"dataVersion"`
	// Match ID
	MatchID string `json:"matchId"`
	// List of participant PUUIDs
	Participants []string `json:"participants"`
}

// MatchInfo contains the data for a specific match
type MatchInfo struct {
	// Unix timestamp for when the game is created on the game server (i.e., the loading screen).
	GameCreation int64 `json:"gameCreation"`
	// Prior to patch 11.20, this field returns the game length in milliseconds calculated
	// from gameEndTimestamp - gameStartTimestamp. Post patch 11.20, this field returns the max
	// timePlayed of any participant in the game in seconds, which makes the behavior of this
	// field consistent with that of match-v4. The best way to handling the change in this field
	// is to treat the value as milliseconds if the gameEndTimestamp field isn't in the response
	// and to treat the value as seconds if gameEndTimestamp is in the response.
	GameDuration int64 `json:"gameDuration"`
	// Unix timestamp for when match ends on the game server. This timestamp can occasionally
	// be significantly longer than when the match "ends". The most reliable way of determining
	// the timestamp for the end of the match would be to add the max time played of any
	// participant to the gameStartTimestamp. This field was added to match-v5 in patch 11.20 on Oct 5th, 2021.
	GameEndTimestamp int64 `json:"gameEndTimestamp"`
	GameID           int64 `json:"gameId"`
	// Please refer to the Game Constants documentation.
	GameMode string `json:"gameMode"`
	GameName string `json:"gameName"`
	// Unix timestamp for when match starts on the game server.
	GameStartTimestamp int64 `json:"gameStartTimestamp"`
	// Please refer to the Game Constants documentation.
	GameType string `json:"gameType"`
	// The first two parts can be used to determine the patch a game was played on.
	GameVersion string `json:"gameVersion"`
	// Please refer to the Game Constants documentation.
	MapID int64 `json:"mapId"`
	// Participant information.
	Participants []*Participant `json:"participants"`
	// Platform where the match was played.
	PlatformID string `json:"platformId"`
	// Please refer to the Game Constants documentation.
	QueueID int64 `json:"queueId"`
	// Team information.
	Teams []*Team `json:"teams"`
	// Tournament code used to generate the match. This field was added to match-v5 in patch 11.13 on June 23rd, 2021.
	TournamentCode string `json:"tournamentCode"`
}

// GetQueue returns the queue this match was played in
func (m *MatchInfo) GetQueue(client *static.Client) (static.Queue, error) {
	return client.GetQueue(m.QueueID)
}

// GetMap returns the map this match was played on
func (m *MatchInfo) GetMap(client *static.Client) (static.Map, error) {
	return client.GetMap(m.MapID)
}

// GetGameType returns the gameType this match was played in
func (m *MatchInfo) GetGameType(client *static.Client) (static.GameType, error) {
	return client.GetGameType(m.GameType)
}

// GetGameMode returns the gameMode this match was played in
func (m *MatchInfo) GetGameMode(client *static.Client) (static.GameMode, error) {
	return client.GetGameMode(m.GameMode)
}

// StatPerks hold stats for a perk
type StatPerks struct {
	Defense int64 `json:"defense"`
	Flex    int64 `json:"flex"`
	Offense int64 `json:"offense"`
}

// Selections contains information about perk selections
type Selections struct {
	Perk int64 `json:"perk"`
	Var1 int64 `json:"var1"`
	Var2 int64 `json:"var2"`
	Var3 int64 `json:"var3"`
}

// Styles holds perk style information
type Styles struct {
	Description string       `json:"description"`
	Selections  []Selections `json:"selections"`
	Style       int64        `json:"style"`
}

// ParticipantPerks holds the perks for a participant in a match
type ParticipantPerks struct {
	StatPerks *StatPerks `json:"statPerks"`
	Styles    []Styles   `json:"styles"`
}

// Participant hold information for a participant of a match
type Participant struct {
	Assists         int64 `json:"assists"`
	BaronKills      int64 `json:"baronKills"`
	BountyLevel     int64 `json:"bountyLevel"`
	ChampExperience int64 `json:"champExperience"`
	ChampLevel      int64 `json:"champLevel"`
	// Prior to patch 11.4, on Feb 18th, 2021, this field returned invalid championIds.
	// We recommend determining the champion based on the championName field for matches played prior to patch 11.4.
	ChampionID   int64  `json:"championId"`
	ChampionName string `json:"championName"`
	// This field is currently only utilized for Kayn's transformations.
	// (Legal values: 0 - None, 1 - Slayer, 2 - Assassin)
	ChampionTransform         int64 `json:"championTransform"`
	ConsumablesPurchased      int64 `json:"consumablesPurchased"`
	DamageDealtToBuildings    int64 `json:"damageDealtToBuildings"`
	DamageDealtToObjectives   int64 `json:"damageDealtToObjectives"`
	DamageDealtToTurrets      int64 `json:"damageDealtToTurrets"`
	DamageSelfMitigated       int64 `json:"damageSelfMitigated"`
	Deaths                    int64 `json:"deaths"`
	DetectorWardsPlaced       int64 `json:"detectorWardsPlaced"`
	DoubleKills               int64 `json:"doubleKills"`
	DragonKills               int64 `json:"dragonKills"`
	FirstBloodAssist          bool  `json:"firstBloodAssist"`
	FirstBloodKill            bool  `json:"firstBloodKill"`
	FirstTowerAssist          bool  `json:"firstTowerAssist"`
	FirstTowerKill            bool  `json:"firstTowerKill"`
	GameEndedInEarlySurrender bool  `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender      bool  `json:"gameEndedInSurrender"`
	GoldEarned                int64 `json:"goldEarned"`
	GoldSpent                 int64 `json:"goldSpent"`
	// Both individualPosition and teamPosition are computed by the game server and are
	// different versions of the most likely position played by a player. The individualPosition
	// is the best guess for which position the player actually played in isolation of
	// anything else. The teamPosition is the best guess for which position the player
	// actually played if we add the constraint that each team must have one top player, one
	// jungle, one middle, etc. Generally the recommendation is to use the teamPosition field
	// over the individualPosition field.
	IndividualPosition             string            `json:"individualPosition"`
	InhibitorKills                 int64             `json:"inhibitorKills"`
	InhibitorTakedowns             int64             `json:"inhibitorTakedowns"`
	InhibitorsLost                 int64             `json:"inhibitorsLost"`
	Item0                          int64             `json:"item0"`
	Item1                          int64             `json:"item1"`
	Item2                          int64             `json:"item2"`
	Item3                          int64             `json:"item3"`
	Item4                          int64             `json:"item4"`
	Item5                          int64             `json:"item5"`
	Item6                          int64             `json:"item6"`
	ItemsPurchased                 int64             `json:"itemsPurchased"`
	KillingSprees                  int64             `json:"killingSprees"`
	Kills                          int64             `json:"kills"`
	Lane                           string            `json:"lane"`
	LargestCriticalStrike          int64             `json:"largestCriticalStrike"`
	LargestKillingSpree            int64             `json:"largestKillingSpree"`
	LargestMultiKill               int64             `json:"largestMultiKill"`
	LongestTimeSpentLiving         int64             `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int64             `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int64             `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int64             `json:"magicDamageTaken"`
	NeutralMinionsKilled           int64             `json:"neutralMinionsKilled"`
	NexusKills                     int64             `json:"nexusKills"`
	NexusLost                      int64             `json:"nexusLost"`
	NexusTakedowns                 int64             `json:"nexusTakedowns"`
	ObjectivesStolen               int64             `json:"objectivesStolen"`
	ObjectivesStolenAssists        int64             `json:"objectivesStolenAssists"`
	ParticipantID                  int64             `json:"participantId"`
	PentaKills                     int64             `json:"pentaKills"`
	Perks                          *ParticipantPerks `json:"perks"`
	PhysicalDamageDealt            int64             `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int64             `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken            int64             `json:"physicalDamageTaken"`
	ProfileIcon                    int64             `json:"profileIcon"`
	PUUID                          string            `json:"puuid"`
	QuadraKills                    int64             `json:"quadraKills"`
	RiotIDGameName                 string            `json:"riotIdGameName"`
	RiotIDName                     string            `json:"riotIdName"`
	RiotIDTagline                  string            `json:"riotIdTagline"`
	Role                           string            `json:"role"`
	SightWardsBoughtInGame         int64             `json:"sightWardsBoughtInGame"`
	Spell1Casts                    int64             `json:"spell1Casts"`
	Spell2Casts                    int64             `json:"spell2Casts"`
	Spell3Casts                    int64             `json:"spell3Casts"`
	Spell4Casts                    int64             `json:"spell4Casts"`
	Summoner1Casts                 int64             `json:"summoner1Casts"`
	Summoner1ID                    int64             `json:"summoner1Id"`
	Summoner2Casts                 int64             `json:"summoner2Casts"`
	Summoner2ID                    int64             `json:"summoner2Id"`
	SummonerID                     string            `json:"summonerId"`
	SummonerLevel                  int64             `json:"summonerLevel"`
	SummonerName                   string            `json:"summonerName"`
	TeamEarlySurrendered           bool              `json:"teamEarlySurrendered"`
	TeamID                         int64             `json:"teamId"`
	// Both individualPosition and teamPosition are computed by the game server and are
	// different versions of the most likely position played by a player. The individualPosition
	// is the best guess for which position the player actually played in isolation of
	// anything else. The teamPosition is the best guess for which position the player
	// actually played if we add the constraint that each team must have one top player, one
	// jungle, one middle, etc. Generally the recommendation is to use the teamPosition field
	// over the individualPosition field.
	TeamPosition                   string `json:"teamPosition"`
	TimeCCingOthers                int64  `json:"timeCCingOthers"`
	TimePlayed                     int64  `json:"timePlayed"`
	TotalDamageDealt               int64  `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int64  `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int64  `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int64  `json:"totalDamageTaken"`
	TotalHeal                      int64  `json:"totalHeal"`
	TotalHealsOnTeammates          int64  `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int64  `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int64  `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int64  `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int64  `json:"totalUnitsHealed"`
	TripleKills                    int64  `json:"tripleKills"`
	TrueDamageDealt                int64  `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int64  `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int64  `json:"trueDamageTaken"`
	TurretKills                    int64  `json:"turretKills"`
	TurretTakedowns                int64  `json:"turretTakedowns"`
	TurretsLost                    int64  `json:"turretsLost"`
	UnrealKills                    int64  `json:"unrealKills"`
	VisionScore                    int64  `json:"visionScore"`
	VisionWardsBoughtInGame        int64  `json:"visionWardsBoughtInGame"`
	WardsKilled                    int64  `json:"wardsKilled"`
	WardsPlaced                    int64  `json:"wardsPlaced"`
	Win                            bool   `json:"win"`
}

// GetSummoner returns the summoner info for this player
func (p *Participant) GetSummoner(client *Client) (*Summoner, error) {
	return client.Summoner.GetByPUUID(p.PUUID)
}

// GetProfileIcon returns the profile icon data for this player
func (p *Participant) GetProfileIcon(client *datadragon.Client) (datadragon.ProfileIcon, error) {
	return client.GetProfileIcon(p.ProfileIcon)
}

// GetChampion returns the champion played by this participant
func (p *Participant) GetChampion(client *datadragon.Client) (datadragon.ChampionDataExtended, error) {
	return client.GetChampionByID(strconv.FormatInt(p.ChampionID, 10))
}

// GetSpell1 returns the first summoner spell of this participant
func (p *Participant) GetSpell1(client *datadragon.Client) (datadragon.SummonerSpell, error) {
	return client.GetSummonerSpell(strconv.FormatInt(p.Summoner1ID, 10))
}

// GetSpell2 returns the second summoner spell of this participant
func (p *Participant) GetSpell2(client *datadragon.Client) (datadragon.SummonerSpell, error) {
	return client.GetSummonerSpell(strconv.FormatInt(p.Summoner2ID, 10))
}

// GetItem0 returns the item in slot 0
func (p *Participant) GetItem0(client *datadragon.Client) (datadragon.Item, error) {
	return client.GetItem(strconv.FormatInt(p.Item0, 10))
}

// GetItem1 returns the item in slot 1
func (p *Participant) GetItem1(client *datadragon.Client) (datadragon.Item, error) {
	return client.GetItem(strconv.FormatInt(p.Item1, 10))
}

// GetItem2 returns the item in slot 2
func (p *Participant) GetItem2(client *datadragon.Client) (datadragon.Item, error) {
	return client.GetItem(strconv.FormatInt(p.Item2, 10))
}

// GetItem3 returns the item in slot 3
func (p *Participant) GetItem3(client *datadragon.Client) (datadragon.Item, error) {
	return client.GetItem(strconv.FormatInt(p.Item3, 10))
}

// GetItem4 returns the item in slot 4
func (p *Participant) GetItem4(client *datadragon.Client) (datadragon.Item, error) {
	return client.GetItem(strconv.FormatInt(p.Item4, 10))
}

// GetItem5 returns the item in slot 5
func (p *Participant) GetItem5(client *datadragon.Client) (datadragon.Item, error) {
	return client.GetItem(strconv.FormatInt(p.Item5, 10))
}

// GetItem6 returns the item in slot 6
func (p *Participant) GetItem6(client *datadragon.Client) (datadragon.Item, error) {
	return client.GetItem(strconv.FormatInt(p.Item6, 10))
}

// TeamBan is a champion banned by a team
type TeamBan struct {
	// Turn during which the champion was banned.
	PickTurn int64 `json:"pickTurn"`
	// Banned championId.
	ChampionID int64 `json:"championId"`
}

// GetChampion returns the champion that was banned
func (b *TeamBan) GetChampion(client *datadragon.Client) (datadragon.ChampionDataExtended, error) {
	return client.GetChampionByID(strconv.FormatInt(b.ChampionID, 10))
}

// Objective holds information for a single objective
type Objective struct {
	First bool  `json:"first"`
	Kills int64 `json:"kills"`
}

// Objectives holds info for a teeam's objeectives
type Objectives struct {
	Baron      Objective `json:"baron"`
	Champion   Objective `json:"champion"`
	Dragon     Objective `json:"dragon"`
	Inhibitor  Objective `json:"inhibitor"`
	RiftHerald Objective `json:"riftHerald"`
	Tower      Objective `json:"tower"`
}

// Team holds information for a team in a match
type Team struct {
	Bans       []*TeamBan `json:"bans"`
	Objectives Objectives `json:"objectives"`
	TeamID     int64      `json:"teamId"`
	Win        bool       `json:"win"`
}

// MatchTimeline contains timeline frames for a match
type MatchTimeline struct {
	Frames   []*MatchFrame `json:"frames"`
	Interval int64         `json:"frameInterval"`
}

// MatchFrame is a single frame in the timeline of a game
type MatchFrame struct {
	Timestamp         int64                        `json:"timestamp"`
	ParticipantFrames map[string]*ParticipantFrame `json:"participantFrames"`
	Events            []*MatchEvent                `json:"events"`
}

// ParticipantFrame contains information about a participant in a game at a single timestamp
type ParticipantFrame struct {
	TotalGold           int64          `json:"totalGold"`
	TeamScore           int64          `json:"teamScore"`
	ParticipantID       int64          `json:"participantId"`
	Level               int64          `json:"level"`
	CurrentGold         int64          `json:"currentGold"`
	MinionsKilled       int64          `json:"minionsKilled"`
	DominionScore       int64          `json:"dominionScore"`
	Position            *MatchPosition `json:"position"`
	XP                  int64          `json:"xp"`
	JungleMinionsKilled int64          `json:"jungleMinionsKilled"`
}

// MatchEventType is the type of an event
type MatchEventType string

// All legal value for match event types
const (
	MatchEventTypeChampionKill     MatchEventType = "CHAMPION_KILL"
	MatchEventTypeWardPlaced                      = "WARD_PLACED"
	MatchEventTypeWardKill                        = "WARD_KILL"
	MatchEventTypeBuildingKill                    = "BUILDING_KILL"
	MatchEventTypeEliteMonsterKill                = "ELITE_MONSTER_KILL"
	MatchEventTypeItemPurchased                   = "ITEM_PURCHASED"
	MatchEventTypeItemSold                        = "ITEM_SOLD"
	MatchEventTypeItemDestroyed                   = "ITEM_DESTROYED"
	MatchEventTypeItemUndo                        = "ITEM_UNDO"
	MatchEventTypeSkillLevelUp                    = "SKILL_LEVEL_UP"
	MatchEventTypeAscendedEvent                   = "ASCENDED_EVENT"
	MatchEventTypeCapturePoint                    = "CAPTURE_POINT"
	MatchEventTypePoroKingSummon                  = "PORO_KING_SUMMON"
)

var (
	// MatchEventTypes is a list of all available match events
	MatchEventTypes = []MatchEventType{
		MatchEventTypeChampionKill,
		MatchEventTypeWardPlaced,
		MatchEventTypeWardKill,
		MatchEventTypeBuildingKill,
		MatchEventTypeEliteMonsterKill,
		MatchEventTypeItemPurchased,
		MatchEventTypeItemSold,
		MatchEventTypeItemDestroyed,
		MatchEventTypeItemUndo,
		MatchEventTypeSkillLevelUp,
		MatchEventTypeAscendedEvent,
		MatchEventTypeCapturePoint,
		MatchEventTypePoroKingSummon,
	}
)

// MatchEvent is an event in a match at a certain timestamp
type MatchEvent struct {
	EventType               string          `json:"eventType"`
	TowerType               string          `json:"towerType"`
	TeamID                  int64           `json:"teamId"`
	AscendedType            string          `json:"ascendedType"`
	KillerID                int64           `json:"killerId"`
	LevelUpType             string          `json:"levelUpType"`
	PointCaptured           string          `json:"pointCaptured"`
	AssistingParticipantIDs []int64         `json:"assistingParticipantIds"`
	WardType                string          `json:"wardType"`
	MonsterType             string          `json:"monsterType"`
	Type                    *MatchEventType `json:"type"`
	SkillSlot               int64           `json:"skillSlot"`
	VictimID                int64           `json:"victimId"`
	Timestamp               int64           `json:"timestamp"`
	AfterID                 int64           `json:"afterId"`
	MonsterSubType          string          `json:"monsterSubType"`
	LaneType                string          `json:"laneType"`
	ItemID                  int64           `json:"itemId"`
	ParticipantID           int64           `json:"participantId"`
	BuildingType            string          `json:"buildingType"`
	CreatorID               int64           `json:"creatorId"`
	Position                *MatchPosition  `json:"position"`
	BeforeID                int64           `json:"beforeId"`
}

// GetItem returns the item for this event
func (e *MatchEvent) GetItem(client *datadragon.Client) (datadragon.Item, error) {
	return client.GetItem(strconv.FormatInt(e.ItemID, 10))
}

// MatchPosition is a position on the map in a game
type MatchPosition struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

// GameInfo contains information about an ongoing game
type GameInfo struct {
	GameID            int64                     `json:"gameId"`
	GameStartTime     int64                     `json:"gameStartTime"`
	PlatformID        string                    `json:"platformId"`
	GameMode          string                    `json:"gameMode"`
	MapID             int64                     `json:"mapId"`
	GameType          string                    `json:"gameType"`
	BannedChampions   []*BannedChampion         `json:"bannedChampions"`
	Observers         *Observer                 `json:"observers"`
	Participants      []*CurrentGameParticipant `json:"participants"`
	GameLength        int64                     `json:"gameLength"`
	GameQueueConfigID int64                     `json:"gameQueueConfigId"`
}

// GetMatch returns information about the finished match
func (i *GameInfo) GetMatch(client *Client) (*Match, error) {
	return client.Match.Get(fmt.Sprintf("%v_%v", strings.ToUpper(string(client.Match.c.Region)), i.GameID))
}

// BannedChampion represents a champion ban during pack/ban phase
type BannedChampion struct {
	PickTurn   int64 `json:"pickTurn"`
	ChampionID int64 `json:"championId"`
	TeamID     int64 `json:"teamId"`
}

// GetChampion returns the banned champion
func (c *BannedChampion) GetChampion(client *datadragon.Client) (datadragon.ChampionDataExtended, error) {
	return client.GetChampionByID(strconv.FormatInt(c.ChampionID, 10))
}

// Observer is an observer of an ongoing game
type Observer struct {
	EncryptionKey string `json:"encryptionKey"`
}

// CurrentGameParticipant represents a player in an ongoing game
type CurrentGameParticipant struct {
	ProfileIconID            int64                      `json:"profileIconId"`
	ChampionID               int64                      `json:"championId"`
	SummonerName             string                     `json:"summonerName"`
	GameCustomizationObjects []*GameCustomizationObject `json:"gameCustomizationObjects"`
	Bot                      bool                       `json:"bot"`
	Perks                    *Perks                     `json:"perks"`
	Spell2ID                 int64                      `json:"spell2Id"`
	Spell1ID                 int64                      `json:"spell1Id"`
	TeamID                   int64                      `json:"teamId"`
	SummonerID               string                     `json:"summonerId"`
	PUUID                    string                     `json:"puuid"`
	RiotID                   string                     `json:"riotId"`
}

// GetChampion returns the champion played by this participant
func (p *CurrentGameParticipant) GetChampion(client *datadragon.Client) (datadragon.ChampionDataExtended, error) {
	return client.GetChampionByID(strconv.FormatInt(p.ChampionID, 10))
}

// GetSpell1 returns the first summoner spell of this participant
func (p *CurrentGameParticipant) GetSpell1(client *datadragon.Client) (datadragon.SummonerSpell, error) {
	return client.GetSummonerSpell(strconv.FormatInt(p.Spell1ID, 10))
}

// GetSpell2 returns the second summoner spell of this participant
func (p *CurrentGameParticipant) GetSpell2(client *datadragon.Client) (datadragon.SummonerSpell, error) {
	return client.GetSummonerSpell(strconv.FormatInt(p.Spell2ID, 10))
}

// GameCustomizationObject contains information specific to an ongoing game
type GameCustomizationObject struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

// Perks represents the runes for a player in an ongoing game
type Perks struct {
	PerkStyle    int64   `json:"perkStyle"`
	PerksIDs     []int64 `json:"perkIds"`
	PerkSubStyle int64   `json:"perkSubStyle"`
}

// FeaturedGames represents a list of featured games
type FeaturedGames struct {
	ClientRefreshInterval int64       `json:"clientRefreshInterval"`
	GameList              []*GameInfo `json:"gameList"`
}

// Status contains information about all services in a certain region
type Status struct {
	Name      string     `json:"name"`
	RegionTag string     `json:"region_tag"`
	Hostname  string     `json:"hostname"`
	Services  []*Service `json:"services"`
	Slug      string     `json:"slug"`
	Locales   []string   `json:"locales"`
}

// Service is a service provided by Riot with its status
type Service struct {
	Status    string      `json:"status"`
	Incidents []*Incident `json:"incidents"`
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
}

// Incident contains information about an incident
type Incident struct {
	Active    bool             `json:"active"`
	CreatedAt string           `json:"created_at"`
	ID        int64            `json:"id"`
	Updates   []*StatusMessage `json:"updates"`
}

// StatusMessage contains information about a status message
type StatusMessage struct {
	Severity     string               `json:"severity"`
	Author       string               `json:"author"`
	CreatedAt    string               `json:"created_at"`
	Translations []*StatusTranslation `json:"translations"`
	UpdatedAt    string               `json:"updated_at"`
	Content      string               `json:"content"`
	ID           string               `json:"id"`
}

// StatusTranslation contains the status message content in a certain language
type StatusTranslation struct {
	Locale    string `json:"locale"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at"`
}

// Summoner represents a summoner with several related IDs
type Summoner struct {
	ProfileIconID int64  `json:"profileIconId"`
	Name          string `json:"name"`
	PUUID         string `json:"puuid"`
	SummonerLevel int64  `json:"summonerLevel"`
	RevisionDate  int64  `json:"revisionDate"`
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
}

// LobbyEventList is a wrapper for a list of lobby events in a tournament
type LobbyEventList struct {
	EventList []*LobbyEvent `json:"eventList"`
}

// LobbyEvent represents an event that happened in a tournament lobby
type LobbyEvent struct {
	EventType  string `json:"eventType"`
	SummonerID string `json:"summonerId"`
	Timestamp  string `json:"timestamp"`
}

// Tournament contains the settings of a previously created tournament
type Tournament struct {
	Map          string   `json:"map"`
	Code         string   `json:"code"`
	Spectators   string   `json:"spectators"`
	Region       string   `json:"region"`
	ProviderID   int64    `json:"providerId"`
	TeamSize     int64    `json:"teamSize"`
	Participants []string `json:"participants"`
	PickType     string   `json:"pickType"`
	TournamentID int64    `json:"tournamentId"`
	LobbyName    string   `json:"lobbyName"`
	Password     string   `json:"password"`
	ID           int64    `json:"id"`
	MetaData     string   `json:"metaData"`
}

// TournamentCodeParameters parameters needed to create tournament codes
type TournamentCodeParameters struct {
	// The spectator type of the game. (Legal values: NONE, LOBBYONLY, ALL)
	SpectatorType string `json:"spectatorType"`
	// The team size of the game. Valid values are 1-5.
	TeamSize int64 `json:"teamSize"`
	// The pick type of the game. (Legal values: BLIND_PICK, DRAFT_MODE, ALL_RANDOM, TOURNAMENT_DRAFT)
	PickType string `json:"pickType"`
	// Optional list of encrypted summonerIds in order to validate the players eligible to join the lobby.
	// NOTE: We currently do not enforce participants at the team level, but rather the aggregate of teamOne and
	// teamTwo. We may add the ability to enforce at the team level in the future.
	AllowedSummonerIDs []string `json:"allowedSummonerIds,omitempty"`
	// The map type of the game. (Legal values: SUMMONERS_RIFT, TWISTED_TREELINE, HOWLING_ABYSS)
	MapType string `json:"mapType"`
	// Optional string that may contain any data in any format, if specified at all. Used to denote any custom
	// information about the game.
	Metadata string `json:"metadata"`
}

// TournamentUpdateParameters parameters needed to update an existing tournament
type TournamentUpdateParameters struct {
	// The spectator type (Legal values: NONE, LOBBYONLY, ALL)
	SpectatorType string `json:"spectatorType"`
	// The pick type (Legal values: BLIND_PICK, DRAFT_MODE, ALL_RANDOM, TOURNAMENT_DRAFT)
	PickType string `json:"pickType"`
	// Optional list of encrypted summonerIds in order to validate the players eligible to join the lobby.
	// NOTE: Participants are not enforced at the team level, but rather the aggregate of teamOne and teamTwo.
	AllowedSummonerIDs []string `json:"allowedSummonerIds"`
	// The map type (Legal values: SUMMONERS_RIFT, TWISTED_TREELINE, HOWLING_ABYSS)
	MapType string `json:"mapType"`
}

// TournamentRegistrationParameters parameters required for creating a tournament
type TournamentRegistrationParameters struct {
	// The provider ID to specify the regional registered provider data to associate this tournament.
	ProviderID int64 `json:"providerId"`
	// The optional name of the tournament.
	Name string `json:"name"`
}

// ProviderRegistrationParameters parameters required for registering a provider with tournaments for a region
type ProviderRegistrationParameters struct {
	// The provider's callback URL to which tournament game results in this region should be posted. The URL must be
	// well-formed, use the http or https protocol, and use the default port for the protocol (http URLs must use port
	// 80, https URLs must use port 443).
	URL string `json:"url"`
	// The region in which the provider will be running tournaments.
	// (Legal values: BR, EUNE, EUW, JP, LAN, LAS, NA, OCE, PBE, RU, TR)
	Region string `json:"region"`
}

// ChallengeConfigInfo represents basic challenge configuration information
type ChallengeConfigInfo struct {
	ID             int64                        `json:"id"`
	LocalizedNames map[string]map[string]string `json:"localizedNames"`
	State          string                       `json:"state"`
	Tracking       string                       `json:"tracking"`
	StartTimeStamp int64                        `json:"startTimeStamp"`
	EndTimeStamp   int64                        `json:"endTimeStamp"`
	Leaderboard    bool                         `json:"leaderboard"`
	Thresholds     map[string]float64           `json:"thresholds"`
}

// ChallengePoints contains the settings of a previously created tournament
type ChallengePoints struct {
	Level      string  `json:"level"`
	Current    float32 `json:"current"`
	Max        int32   `json:"max"`
	Percentile float32 `json:"percentile"`
}

// ChallengeInfo represents each challenge info for a player
type ChallengeInfo struct {
	ChallengeID  int32   `json:"challengeid"`
	Percentile   float32 `json:"percentile"`
	Level        string  `json:"level"`
	Value        float32 `json:"value"`
	AchievedTime int64   `json:"achievedtime"`
}

// PlayerClientPreferences holds player preferences
type PlayerClientPreferences struct {
	BannerAccent string  `json:"banneraccent"`
	Title        string  `json:"title"`
	ChallengeID  []int32 `json:"challengeids"`
}

// PlayerInfo contains player information with list of all progressed challenges
type PlayerInfo struct {
	TotalPoints    *ChallengePoints           `json:"totalpoints"`
	CategoryPoints map[string]ChallengePoints `json:"categorypoints"`
	Challenges     []*ChallengeInfo           `json:"challenges"`
	Preferences    *PlayerClientPreferences   `json:"preferences"`
}

// ApexPlayerInfo holds information of top players for each level
type ApexPlayerInfo struct {
	PuuID    string  `json:"puuid"`
	Value    float64 `json:"value"`
	Position int32   `json:"position"`
}
