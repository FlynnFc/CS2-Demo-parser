package main

import "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"

type FinalScore struct {
	CT int
	T  int
}

type PlayerStats struct {
	ID             uint64
	Name           string  `json:"name"`
	TeamName       string  `json:"teamName"`
	Rounds         int     `json:"rounds"`
	Matches        int     `json:"matches"`
	Kills          int     `json:"kills"`
	Assists        int     `json:"assists"`
	Deaths         int     `json:"deaths"`
	Damage         int     `json:"damage"`
	KASTRounds     int     `json:"KAST"`
	MischiefRating float64 `json:"mischiefRating"`
	DetailedStats  DetailedPlayerStats
	RoundEvent     []RoundEvent `json:"roundEvent"`
}

type RoundEvent struct {
	RoundWon      bool
	RoundNumber   int8
	Kills         int8
	Died          bool
	Assists       int8
	Damage        int
	TradeAttempts int8
	TradeKills    int8
	DeathTraded   bool
	EntryKill     bool
	EntryDeath    bool
	Side          common.Team
	Clutch        Clutch
}

type BombStats struct {
	BombsPlanted       int
	BombsDefused       int
	BombDefuseAttempts int
	BombPlantAttempts  int
}

type KillStats struct {
	Wallbangs     int
	Headshots     int
	Blind         int
	NoScope       int
	ThroughSmoke  int
	AirborneKills int
	Entry         EntryStats
	TradeStats    TradeStats
}

type TradeStats struct {
	TotalAttempts  int
	CTTradeKills   int
	CTFailedTrades int
	CTTradedDeaths int
	TTradeKills    int
	TFailedTrades  int
	TTradedDeaths  int
}
type EntryStats struct {
	TotalAttempts int
	CTEntryKills  int
	CTEntryDeaths int
	TEntryKills   int
	TEntryDeaths  int
}

type Clutch struct {
	Type  int8
	Won   bool
	Round int
	Kills int
}

type DetailedPlayerStats struct {
	CTRounds    int
	TRounds     int
	RefundTotal int
	BombStats   BombStats
	KillStats   KillStats
	ClutchStats []Clutch
}

type ResponsePayload struct {
	Payload struct {
		DownloadURL string `json:"download_url"`
	} `json:"payload"`
}
type DownloadRequest struct {
	ResourceURL string `json:"resource_url"`
}

type KillState struct {
	Time     int
	Killer   *common.Player
	Victim   *common.Player
	Assister *common.Player
	TSide    []*common.Player
	CTSide   []*common.Player
}
