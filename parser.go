package main

import (
	"errors"
	"log"
	"os"

	demoinfocs "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func demoParsing(path string) ([]PlayerStats, error) {
	f, err := os.Open(path)
	checkError(err)
	defer f.Close()
	var live = 0
	p := demoinfocs.NewParser(f)
	defer p.Close()
	startingOutput := []PlayerStats{}
	output := []PlayerStats{}

	roundKills := []KillState{}
	var allPlayers []*common.Player
	//Detects start of game
	p.RegisterEventHandler(func(e events.RoundStart) {
		pistolRound := true

		for _, player := range p.GameState().Participants().Playing() {
			if player.Money() != 800 {
				pistolRound = false
			}
		}

		if pistolRound {
			live++
			if live == 1 {
				// Match starts
				startingOutput = []PlayerStats{}
				for _, player := range p.GameState().Participants().Playing() {
					allPlayers = append(allPlayers, player)
					startingOutput = append(startingOutput, PlayerStats{ID: player.SteamID64, Name: player.Name, Rounds: 0, DetailedStats: DetailedPlayerStats{CTRounds: 0, TRounds: 0}})
				}
			}
		}
	})

	p.RegisterEventHandler(func(e events.RoundEnd) {
		//Counting T and CT Rounds per player (Not working on challengerMode but fuck em)
		for i := 0; i < len(startingOutput); i++ {
			for j := 0; j < len(startingOutput); j++ {
				if startingOutput[i].ID == allPlayers[j].SteamID64 {
					if allPlayers[j].Team == 2 {
						startingOutput[i].DetailedStats.TRounds++
					} else if allPlayers[j].Team == 3 {
						startingOutput[i].DetailedStats.CTRounds++
					}
					break
				}
			}
		}
		// VALID ROUND
		if live > 0 && len(roundKills) > 0 {
			rounds := p.GameState().TeamCounterTerrorists().Score() + p.GameState().TeamTerrorists().Score()

			// Initalise round
			for i := 0; i < len(startingOutput); i++ {
				playerRef := findPlayer(startingOutput[i].ID, p.GameState())
				var roundOutcome bool
				if playerRef.Team == e.Winner {
					roundOutcome = true
				} else {
					roundOutcome = false
				}

				startingOutput[i].RoundEvent = append(startingOutput[i].RoundEvent, RoundEvent{RoundNumber: int8(rounds), Died: false, RoundWon: roundOutcome, Side: playerRef.Team})
			}
			// Clutches
			err := determineClutchOutcome(roundKills, p.GameState(), e.Winner, startingOutput)
			if err != nil {
				log.Printf("Failed to parse clutches for round %d\n", rounds)

			}
			// Trade kills, Kills, Assists
			err = determineTradeKills(roundKills, rounds, startingOutput)
			if err != nil {
				log.Printf("Failed to parse Trades for round %d\n", rounds)
				log.Println(err)
			}
			// TODO
			err = determinePlayerSurvival(roundKills, rounds, startingOutput)
			if err != nil {
				log.Printf("Failed to parse surviver checks %d\n", rounds)
			}
			err = determinePlayerKillOrAssist(roundKills, rounds, startingOutput)
			if err != nil {
				log.Printf("Failed to parse Kill and Assist checks %d\n", rounds)
			}
			err = determineEntryKill(roundKills, rounds, startingOutput)
			if err != nil {
				log.Printf("Failed to parse Entry Checks %d\n", rounds)
			}
		}
		roundKills = roundKills[:0]
	})

	p.RegisterEventHandler(func(e events.Kill) {
		// kills = append(kills, e)
		if live > 0 && e.Victim != nil {
			var targetPlayer uint64
			kills = append(kills, e)
			victimTeam := e.Victim.TeamState.Members()
			if e.Killer != nil {
				targetPlayer = e.Killer.SteamID64
				TSide := p.GameState().TeamTerrorists().Members()
				CTSide := p.GameState().TeamCounterTerrorists().Members()
				// Time in seconds since demo started
				time := p.CurrentFrame() / int(p.TickRate())
				roundKills = append(roundKills, KillState{Killer: e.Killer, Victim: e.Victim, Assister: e.Assister, TSide: TSide, CTSide: CTSide, Time: time})
			}

			for i, player := range startingOutput {
				//Player found in group
				if player.ID == targetPlayer {
					//Wallbang Kills
					if e.IsWallBang() {
						startingOutput[i].DetailedStats.KillStats.Wallbangs++
					}
					//Headshot Kills
					if e.IsHeadshot {
						startingOutput[i].DetailedStats.KillStats.Headshots++
					}
					//Blind kills
					if e.AttackerBlind {
						startingOutput[i].DetailedStats.KillStats.Blind++
					}
					if e.ThroughSmoke {
						startingOutput[i].DetailedStats.KillStats.ThroughSmoke++
					}
					if e.Killer.IsAirborne() {
						startingOutput[i].DetailedStats.KillStats.AirborneKills++
					}
					//Entry Kill Stats
					isEntry := true
					for _, player := range victimTeam {
						if !player.IsAlive() {
							isEntry = false
						}
					}

					if isEntry {
						startingOutput[i].DetailedStats.KillStats.Entry.TotalAttempts++
						//2 = T side
						if e.Killer.Team == 2 {
							startingOutput[i].DetailedStats.KillStats.Entry.TEntryKills++
						} else {
							startingOutput[i].DetailedStats.KillStats.Entry.CTEntryKills++
						}
					}
				}
			}
			var deadPlayer uint64
			if e.Victim != nil {
				deadPlayer = e.Victim.SteamID64
			}
			// Die after getting Kill
			deadPlayerIsEntry := true
			for _, player := range victimTeam {
				if !player.IsAlive() {
					deadPlayerIsEntry = false
				}
			}
			// Entry kills Lost
			if deadPlayerIsEntry {
				for i, player := range startingOutput {
					//Player found in group
					if player.ID == deadPlayer {
						startingOutput[i].DetailedStats.KillStats.Entry.TotalAttempts++
						//2 = T side
						if e.Killer.Team == 2 {
							// T First Deaths
							startingOutput[i].DetailedStats.KillStats.Entry.TEntryDeaths++
						} else {
							// CT First Deaths
							startingOutput[i].DetailedStats.KillStats.Entry.CTEntryDeaths++
						}
					}
				}
			}
		}
	})

	//CT and T Rounds
	//Objectives / Round
	//Bomb Plants
	p.RegisterEventHandler(func(e events.BombPlanted) {
		targetPlayer := e.Player.SteamID64
		for i, player := range startingOutput {
			if player.ID == targetPlayer {
				startingOutput[i].DetailedStats.BombStats.BombsPlanted++
				startingOutput[i].DetailedStats.BombStats.BombPlantAttempts++
			}
		}
	})
	//Bomb Defused
	p.RegisterEventHandler(func(e events.BombDefused) {
		targetPlayer := e.Player.SteamID64
		for i, player := range startingOutput {
			if player.ID == targetPlayer {
				startingOutput[i].DetailedStats.BombStats.BombsDefused++
				startingOutput[i].DetailedStats.BombStats.BombDefuseAttempts++
			}
		}
	})
	//Bomb Plant attempts
	p.RegisterEventHandler(func(e events.BombPlantAborted) {
		targetPlayer := e.Player.SteamID64
		for i, player := range startingOutput {
			if player.ID == targetPlayer {
				startingOutput[i].DetailedStats.BombStats.BombPlantAttempts++
			}
		}
	})
	//Bomb Defused attempts
	p.RegisterEventHandler(func(e events.BombDefuseAborted) {
		targetPlayer := e.Player.SteamID64
		for i, player := range startingOutput {
			if player.ID == targetPlayer {
				startingOutput[i].DetailedStats.BombStats.BombDefuseAttempts++
			}
		}
	})

	//TODO
	// Trade Kills
	//TODO
	// Trade Deaths

	p.RegisterEventHandler(func(e events.ItemRefund) {
		player := e.Player.SteamID64
		for i := 0; i < len(startingOutput); i++ {
			if player == startingOutput[i].ID {
				startingOutput[i].DetailedStats.RefundTotal++
			}
		}
	})

	// //Game ended
	p.RegisterEventHandler(func(e events.AnnouncementWinPanelMatch) {
		// allPlayers := p.GameState().Participants().Playing()
		score := FinalScore{p.GameState().TeamCounterTerrorists().Score(), p.GameState().TeamTerrorists().Score()}

		totalRounds := score.CT + score.T

		if score.CT < 12 && score.T < 12 {
			//Not 1 half has been completed so match is likely void
			failed = append(failed, path)
			log.Printf(Red+"Demo %s not valid. Too few rounds player\n"+Reset, path)
			return
		}
		for _, player := range allPlayers {
			stats := playerStatsCalc(player, totalRounds)

			for _, oldPlayer := range startingOutput {

				if stats.ID == oldPlayer.ID {

					KAST := playerKASTCalc(oldPlayer)
					// Create new PlayerStats object
					newPlayerStats := PlayerStats{
						Name:          stats.Name,
						TeamName:      player.TeamState.ClanName(),
						ID:            stats.ID,
						Rounds:        stats.Rounds,
						Matches:       1,
						Kills:         stats.Kills,
						Deaths:        stats.Deaths,
						Damage:        stats.Damage,
						DetailedStats: oldPlayer.DetailedStats,
						RoundEvent:    oldPlayer.RoundEvent,
						KAST:          int8(KAST),
					}

					// Copy oldPlayer's DetailedStats

					// Append newPlayerStats to output
					output = append(output, newPlayerStats)
				}
			}
		}

	})

	// Parse the whole demo
	err = p.ParseToEnd()

	if err != nil {
		log.Println(Red + "Demo ended abruptly! Likely was not a completed game." + Reset)
		log.Println(err)
		failed = append(failed, path)
		return nil, err
	}

	return output, nil
}

func findPlayerIndex(arr1 []PlayerStats, target *common.Player) (int, error) {
	for i, player := range arr1 {
		//Player found in group
		if player.Name == target.Name {
			return i, nil
		}
	}
	return 0, errors.New("ID not found in the slice")
}

func determineClutchOutcome(killEvents []KillState, gameState demoinfocs.GameState, winner common.Team, output []PlayerStats) error {
	alivePlayersT := make(map[*common.Player]bool)
	alivePlayersCT := make(map[*common.Player]bool)
	// Initialize alive players for both teams
	for _, player := range gameState.Team(common.TeamTerrorists).Members() {
		alivePlayersT[player] = true
	}
	for _, player := range gameState.Team(common.TeamCounterTerrorists).Members() {
		alivePlayersCT[player] = true
	}

	potentialClutchWinner := &common.Player{}
	clutchSituation := 0
	oneOnOneSituation := false

	for _, kill := range killEvents {
		if kill.Victim != nil {
			if kill.Victim.Team == common.TeamTerrorists {
				delete(alivePlayersT, kill.Victim)
			} else {
				delete(alivePlayersCT, kill.Victim)
			}
		}

		// Check if either team is in a clutch situation
		if len(alivePlayersT) == 1 && len(alivePlayersCT) > 1 {
			for player := range alivePlayersT {
				potentialClutchWinner = player
			}
			if clutchSituation == 0 {
				clutchSituation = len(alivePlayersCT)
			}

			// fmt.Printf("1v%d - %s\n", clutchSituation, potentialClutchWinner.Name)
		} else if len(alivePlayersCT) == 1 && len(alivePlayersT) > 1 {
			for player := range alivePlayersCT {
				potentialClutchWinner = player
			}
			if clutchSituation == 0 {
				clutchSituation = len(alivePlayersT)
			}
			// fmt.Printf("1v%d - %s\n", clutchSituation, potentialClutchWinner.Name)
		}

		// Detect if it's a 1v1 situation
		if len(alivePlayersT) == 1 && len(alivePlayersCT) == 1 {
			oneOnOneSituation = true
		}
	}

	// Check the final state of the clutch
	if clutchSituation > 0 || oneOnOneSituation {
		lastKill := killEvents[len(killEvents)-1]
		potentiali, err := findPlayerIndex(output, potentialClutchWinner)
		round := gameState.TeamCounterTerrorists().Score() + gameState.TeamTerrorists().Score()
		if lastKill.Killer != nil && potentialClutchWinner.Team == winner && err == nil {
			if !oneOnOneSituation {
				//Apply win on main clutcher
				if potentiali != -1 {
					// 1.0
					output[potentiali].DetailedStats.ClutchStats = append(output[potentiali].DetailedStats.ClutchStats, Clutch{Type: int8(clutchSituation), Won: true, Round: round, Kills: clutchSituation})

					// 2.0
					target := currentRoundIdx(round, output[potentiali].RoundEvent)
					if target != -1 {
						output[potentiali].RoundEvent[target].Clutch.Kills = clutchSituation
						output[potentiali].RoundEvent[target].Clutch.Won = true
						output[potentiali].RoundEvent[target].Clutch.Type = int8(clutchSituation)
					}
					// fmt.Printf("%s won 1v%d\n", potentialClutchWinner.Name, clutchSituation)
				}
				return nil
			} else {
				if potentialClutchWinner.IsAlive() {
					//Apply win on main clutcher

					// 1.0
					output[potentiali].DetailedStats.ClutchStats = append(output[potentiali].DetailedStats.ClutchStats, Clutch{Type: int8(clutchSituation), Won: true, Round: round, Kills: clutchSituation})

					// 2.0
					target := currentRoundIdx(round, output[potentiali].RoundEvent)
					if target != -1 {
						output[potentiali].RoundEvent[target].Clutch.Kills = clutchSituation
						output[potentiali].RoundEvent[target].Clutch.Won = true
						output[potentiali].RoundEvent[target].Clutch.Type = int8(clutchSituation)
					}

					//Set victim to losing a 1v1
					loser, err := findPlayerIndex(output, lastKill.Victim)
					if err == nil {
						// 1.0
						output[loser].DetailedStats.ClutchStats = append(output[loser].DetailedStats.ClutchStats, Clutch{Type: 1, Won: false, Round: round, Kills: 0})
						// fmt.Printf("%s won 1v%d -- %s lost 1v1\n", potentialClutchWinner.Name, clutchSituation, lastKill.Victim.Name)

						// 2.0
						target := currentRoundIdx(round, output[loser].RoundEvent)
						if target != -1 {
							output[loser].RoundEvent[target].Clutch.Kills = 0
							output[loser].RoundEvent[target].Clutch.Won = false
							output[loser].RoundEvent[target].Clutch.Type = 1
						}
					}
				} else {
					// 1.0
					output[potentiali].DetailedStats.ClutchStats = append(output[potentiali].DetailedStats.ClutchStats, Clutch{Type: int8(clutchSituation), Won: true, Round: round, Kills: clutchSituation})

					// 2.0
					winnerTarget := currentRoundIdx(round, output[potentiali].RoundEvent)
					if winnerTarget != -1 {
						output[potentiali].RoundEvent[winnerTarget].Clutch.Kills = 0
						output[potentiali].RoundEvent[winnerTarget].Clutch.Won = true
						output[potentiali].RoundEvent[winnerTarget].Clutch.Type = int8(clutchSituation)
					}
					loser, err := findPlayerIndex(output, lastKill.Killer)
					if err == nil {
						// 1.0
						output[loser].DetailedStats.ClutchStats = append(output[loser].DetailedStats.ClutchStats, Clutch{Type: 1, Won: false, Round: round, Kills: 1})
						// fmt.Printf("%s won 1v%d -- %s lost 1v1\n", potentialClutchWinner.Name, clutchSituation, lastKill.Killer.Name)

						// 2.0
						loserTarget := currentRoundIdx(round, output[loser].RoundEvent)
						if loserTarget != -1 {
							output[loser].RoundEvent[loserTarget].Clutch.Kills = 1
							output[loser].RoundEvent[loserTarget].Clutch.Won = false
							output[loser].RoundEvent[loserTarget].Clutch.Type = 1
						}
					}

				}
				return nil
			}
		} else {
			//Lost the clutch
			if oneOnOneSituation {
				loser, err := findPlayerIndex(output, potentialClutchWinner)

				if err == nil {
					// 1.0
					output[loser].DetailedStats.ClutchStats = append(output[loser].DetailedStats.ClutchStats, Clutch{Type: int8(clutchSituation), Won: false, Round: round, Kills: clutchSituation - 1})

					// 2.0
					loserTarget := currentRoundIdx(round, output[loser].RoundEvent)
					if loserTarget != -1 {
						output[loser].RoundEvent[loserTarget].Clutch.Kills = clutchSituation - 1
						output[loser].RoundEvent[loserTarget].Clutch.Won = false
						output[loser].RoundEvent[loserTarget].Clutch.Type = int8(clutchSituation)
					}
				}
				winner, err := findPlayerIndex(output, lastKill.Killer)
				if err == nil {
					// 1.0
					output[winner].DetailedStats.ClutchStats = append(output[winner].DetailedStats.ClutchStats, Clutch{Type: 1, Won: true, Round: round, Kills: 1})
					// fmt.Printf("%s lost 1v%d -- %s won 1v1\n", potentialClutchWinner.Name, clutchSituation, lastKill.Killer.Name)

					// 2.0
					winnerTarget := currentRoundIdx(round, output[winner].RoundEvent)
					if winnerTarget != -1 {
						output[winner].RoundEvent[winnerTarget].Clutch.Kills = 1
						output[winner].RoundEvent[winnerTarget].Clutch.Won = true
						output[winner].RoundEvent[winnerTarget].Clutch.Type = 1
					}
				}
				return nil
			} else {
				leftAlive := CountAliveOnWinningTeam(killEvents, winner, gameState)
				// 1.0
				if potentiali != -1 {
					output[potentiali].DetailedStats.ClutchStats = append(output[potentiali].DetailedStats.ClutchStats, Clutch{Type: int8(clutchSituation), Won: false, Round: round, Kills: clutchSituation - leftAlive})
				}
				// fmt.Printf("%s lost 1v%d\n", potentialClutchWinner.Name, clutchSituation)

				// 2.0
				target := currentRoundIdx(round, output[potentiali].RoundEvent)
				if target != -1 {
					output[potentiali].RoundEvent[target].Clutch.Kills = clutchSituation
					output[potentiali].RoundEvent[target].Clutch.Won = true
					output[potentiali].RoundEvent[target].Clutch.Type = int8(clutchSituation)
				}

				return nil
			}
		}
	}
	return nil
}

func CountAliveOnWinningTeam(killEvents []KillState, winningTeam common.Team, gameState demoinfocs.GameState) int {
	alivePlayers := make(map[*common.Player]bool)

	// Initialize all players on the winning team as alive
	for _, player := range gameState.Team(winningTeam).Members() {
		alivePlayers[player] = true
	}
	aliveCount := 5
	// Process each kill event
	for _, kill := range killEvents {
		if kill.Victim != nil && kill.Victim.Team == winningTeam {
			aliveCount--
		}
	}

	return aliveCount
}

func playerStatsCalc(player *common.Player, totalRounds int) PlayerStats {
	var output = PlayerStats{ID: player.SteamID64, Name: player.Name, Rounds: totalRounds, Kills: player.Kills(), Deaths: player.Deaths(), Assists: player.Assists(), Damage: player.TotalDamage(), Matches: 1}
	return output
}

func determineTradeKills(killEvents []KillState, round int, output []PlayerStats) error {

outer:
	for i, kill := range killEvents {
		// If a player has killed another player

		// We want to track...

		// If a 2nd player has died to the same person (this is an attempted trade) -
		// Also want to track if a player has been traded after they died (traded deaths)

		if kill.Killer != nil && kill.Victim != nil {

			// Loop through the rest of the kills and see if they have been traded
			for j := i + 1; j < len(killEvents); j++ {
				potentialTrade := killEvents[j]
				if potentialTrade.Killer != nil && potentialTrade.Victim != nil {
					// If kills are more than 6 seconds old they arent relevant
					if potentialTrade.Time-kill.Time > 6 {
						continue outer
					}
					// Player who got a kill was killed
					if potentialTrade.Victim.SteamID64 == kill.Killer.SteamID64 {
						traderIdx, err := findPlayerIndex(output, potentialTrade.Killer)
						if err != nil {
							// log.Printf("player %s not found", kill.Killer.Name)
							continue
						}
						playerTraded, err := findPlayerIndex(output, kill.Victim)
						if err != nil {
							// log.Printf("player %s not found", kill.Victim.Name)
							continue
						}
						// 1.0
						output[traderIdx].DetailedStats.KillStats.TradeStats.TotalAttempts++

						// 2.0
						traderRoundIdx := currentRoundIdx(round, output[traderIdx].RoundEvent)
						tradedPlayerIdx := currentRoundIdx(round, output[playerTraded].RoundEvent)

						// Increments players trade kills. Also increments original victims traded deaths.
						if traderRoundIdx != -1 {
							output[traderIdx].RoundEvent[traderRoundIdx].TradeAttempts++
							output[traderIdx].RoundEvent[traderRoundIdx].TradeKills++
						}
						if tradedPlayerIdx != -1 {
							output[playerTraded].RoundEvent[tradedPlayerIdx].DeathTraded = true
						}

						if potentialTrade.Killer.Team == common.TeamCounterTerrorists {
							// 1.0
							output[traderIdx].DetailedStats.KillStats.TradeStats.CTTradeKills++
							output[playerTraded].DetailedStats.KillStats.TradeStats.CTTradedDeaths++
							continue outer
						} else {
							output[traderIdx].DetailedStats.KillStats.TradeStats.TTradeKills++
							output[playerTraded].DetailedStats.KillStats.TradeStats.TTradedDeaths++
							continue outer
						}
					}
					// Player who got a kill denied a trade attempt
					if potentialTrade.Killer.SteamID64 == kill.Killer.SteamID64 {
						failedTrader, err := findPlayerIndex(output, potentialTrade.Victim)
						if err != nil {
							// log.Printf("player %s not found", kill.Killer.Name)
							continue
						}

						// 2.0
						failedIdx := currentRoundIdx(round, output[failedTrader].RoundEvent)
						if failedIdx != -1 {
							output[failedTrader].RoundEvent[failedIdx].TradeAttempts++
						}

						// 1.0
						output[failedTrader].DetailedStats.KillStats.TradeStats.TotalAttempts++

						if potentialTrade.Killer.Team == common.TeamCounterTerrorists {
							output[failedTrader].DetailedStats.KillStats.TradeStats.CTFailedTrades++
							continue outer
						} else {
							output[failedTrader].DetailedStats.KillStats.TradeStats.TFailedTrades++
							continue outer
						}
					}
				}
			}
		}
	}
	return nil
}

func determinePlayerSurvival(killEvents []KillState, roundNumber int, output []PlayerStats) error {

	for _, kill := range killEvents {
		if kill.Victim != nil {
			idx, err := findPlayerIndex(output, kill.Victim)
			if err != nil {
				continue
			}
			roundIdx := currentRoundIdx(roundNumber, output[idx].RoundEvent)
			if roundIdx != -1 {
				output[idx].RoundEvent[roundIdx].Died = true
			}
		}
	}
	return nil
}
func determinePlayerKillOrAssist(killEvents []KillState, roundNumber int, output []PlayerStats) error {

	for _, kill := range killEvents {
		if kill.Killer != nil {
			idx, err := findPlayerIndex(output, kill.Killer)
			if err != nil {
				continue
			}
			roundIdx := currentRoundIdx(roundNumber, output[idx].RoundEvent)
			if roundIdx != -1 {
				output[idx].RoundEvent[roundIdx].Kills++
			}
		}
		if kill.Assister != nil {
			idx, err := findPlayerIndex(output, kill.Assister)
			if err != nil {
				continue
			}
			roundIdx := currentRoundIdx(roundNumber, output[idx].RoundEvent)
			if roundIdx != -1 {
				output[idx].RoundEvent[roundIdx].Assists++
			}
		}
	}
	return nil
}
func determineEntryKill(killEvents []KillState, roundNumber int, output []PlayerStats) error {
	for _, entry := range killEvents {
		killer, err := findPlayerIndex(output, entry.Killer)
		if err != nil {
			return err
		}
		victim, err := findPlayerIndex(output, entry.Victim)
		if err != nil {
			return err
		}
		killerRoundIdx := currentRoundIdx(roundNumber, output[killer].RoundEvent)
		victimRoundIdx := currentRoundIdx(roundNumber, output[victim].RoundEvent)

		if killerRoundIdx != -1 && victimRoundIdx != -1 && entry.Victim != nil {
			output[killer].RoundEvent[killerRoundIdx].EntryKill = true
			output[victim].RoundEvent[victimRoundIdx].EntryDeath = true
			return nil
		}
	}
	return nil
}
func currentRoundIdx(target int, haystack []RoundEvent) int {
	for i, curr := range haystack {
		if curr.RoundNumber == int8(target) {
			return i
		}
	}
	return -1
}

func playerKASTCalc(player PlayerStats) int {
	rounds := player.DetailedStats.TRounds + player.DetailedStats.CTRounds
	KASTRounds := 0
	for _, round := range player.RoundEvent {
		if round.Kills > 0 || round.Assists > 0 || !round.Died || round.TradeKills > 0 {
			KASTRounds++
		}
	}

	if KASTRounds == 0 {
		return -1
	}
	fullKAST := float64(KASTRounds) / float64(rounds)
	return int(fullKAST * 100)
}

func findPlayer(target uint64, game demoinfocs.GameState) *common.Player {
	for _, player := range game.Participants().All() {
		if player.SteamID64 == target {
			return player
		}
	}
	return nil
}
