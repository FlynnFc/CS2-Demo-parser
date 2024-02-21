package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

func excelExporter(allPlayers []PlayerStats) {
	fmt.Println(Yellow + "Building spreadsheet..." + Reset)

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("rawdata")
	checkError(err)

	//
	//Header row initialising
	//
	headerRow := sheet.AddRow()
	headerRow.AddCell().SetString("ID")
	headerRow.AddCell().SetString("Name")
	headerRow.AddCell().SetString("Team Name")
	headerRow.AddCell().SetString("Matches")
	headerRow.AddCell().SetString("Rounds")
	headerRow.AddCell().SetString("Kills")
	headerRow.AddCell().SetString("Assists")
	headerRow.AddCell().SetString("Deaths")
	headerRow.AddCell().SetString("Damage")
	headerRow.AddCell().SetString("ADR") // Average Damage per Round
	headerRow.AddCell().SetString("KAST")
	headerRow.AddCell().SetString("Mischief Rating")
	headerRow.AddCell().SetString("CT Rounds")
	headerRow.AddCell().SetString("T Rounds")
	headerRow.AddCell().SetString("Refund Total")
	headerRow.AddCell().SetString("Bombs Planted")
	headerRow.AddCell().SetString("Bombs Defused")
	headerRow.AddCell().SetString("Bomb Defuse Attempts")
	headerRow.AddCell().SetString("Bomb Plant Attempts")
	headerRow.AddCell().SetString("Wallbangs")
	headerRow.AddCell().SetString("Headshots")
	headerRow.AddCell().SetString("Blind Kills")
	headerRow.AddCell().SetString("No Scopes")
	headerRow.AddCell().SetString("Through Smoke Kills")
	headerRow.AddCell().SetString("Airborne Kills")
	headerRow.AddCell().SetString("Total Trade Attempts")
	headerRow.AddCell().SetString("CT Trade Kills")
	headerRow.AddCell().SetString("CT Failed Trades")
	headerRow.AddCell().SetString("CT Traded Deaths")
	headerRow.AddCell().SetString("T Trade Kills")
	headerRow.AddCell().SetString("T Failed Trades")
	headerRow.AddCell().SetString("T Traded Deaths")
	headerRow.AddCell().SetString("Total Entry Attempts")
	headerRow.AddCell().SetString("CT Entry Kills")
	headerRow.AddCell().SetString("CT Entry Deaths")
	headerRow.AddCell().SetString("T Entry Kills")
	headerRow.AddCell().SetString("T Entry Deaths")
	for i := 1; i <= 5; i++ {
		headerRow.AddCell().SetString(fmt.Sprintf("1v%d Attempts", i))
		headerRow.AddCell().SetString(fmt.Sprintf("1v%d Wins", i))
	}

	for _, player := range allPlayers {
		var id = strconv.FormatUint(player.ID, 10)
		row := sheet.AddRow()
		row.AddCell().SetString(id)
		row.AddCell().SetString(player.Name)
		row.AddCell().SetString(player.TeamName)
		row.AddCell().SetInt(player.Matches)
		row.AddCell().SetInt(player.Rounds)
		row.AddCell().SetInt(player.Kills)
		row.AddCell().SetInt(player.Assists)
		row.AddCell().SetInt(player.Deaths)
		row.AddCell().SetInt(player.Damage)
		row.AddCell().SetInt(player.Damage / player.Rounds)
		row.AddCell().SetInt(player.Rounds / player.Rounds)
		row.AddCell().SetFloat(player.MischiefRating)
		row.AddCell().SetInt(player.DetailedStats.CTRounds)
		row.AddCell().SetInt(player.DetailedStats.TRounds)
		row.AddCell().SetInt(player.DetailedStats.RefundTotal)
		row.AddCell().SetInt(player.DetailedStats.BombStats.BombsPlanted)
		row.AddCell().SetInt(player.DetailedStats.BombStats.BombsDefused)
		row.AddCell().SetInt(player.DetailedStats.BombStats.BombDefuseAttempts)
		row.AddCell().SetInt(player.DetailedStats.BombStats.BombPlantAttempts)
		row.AddCell().SetInt(player.DetailedStats.KillStats.Wallbangs)
		row.AddCell().SetInt(player.DetailedStats.KillStats.Headshots)
		row.AddCell().SetInt(player.DetailedStats.KillStats.Blind)
		row.AddCell().SetInt(player.DetailedStats.KillStats.NoScope)
		row.AddCell().SetInt(player.DetailedStats.KillStats.ThroughSmoke)
		row.AddCell().SetInt(player.DetailedStats.KillStats.AirborneKills)

		// TradeStats
		row.AddCell().SetInt(player.DetailedStats.KillStats.TradeStats.TotalAttempts)
		row.AddCell().SetInt(player.DetailedStats.KillStats.TradeStats.CTTradeKills)
		row.AddCell().SetInt(player.DetailedStats.KillStats.TradeStats.CTFailedTrades)
		row.AddCell().SetInt(player.DetailedStats.KillStats.TradeStats.CTTradedDeaths)
		row.AddCell().SetInt(player.DetailedStats.KillStats.TradeStats.TTradeKills)
		row.AddCell().SetInt(player.DetailedStats.KillStats.TradeStats.TFailedTrades)
		row.AddCell().SetInt(player.DetailedStats.KillStats.TradeStats.TTradedDeaths)

		// EntryStats
		row.AddCell().SetInt(player.DetailedStats.KillStats.Entry.TotalAttempts)
		row.AddCell().SetInt(player.DetailedStats.KillStats.Entry.CTEntryKills)
		row.AddCell().SetInt(player.DetailedStats.KillStats.Entry.CTEntryDeaths)
		row.AddCell().SetInt(player.DetailedStats.KillStats.Entry.TEntryKills)
		row.AddCell().SetInt(player.DetailedStats.KillStats.Entry.TEntryDeaths)

		clutchAttempts := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
		clutchWins := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}

		// Count attempts and wins for each clutch type
		for _, clutch := range player.DetailedStats.ClutchStats {
			if clutch.Type >= 1 && clutch.Type <= 5 {
				clutchAttempts[int(clutch.Type)]++
				if clutch.Won {
					clutchWins[int(clutch.Type)]++
				}
			}
		}

		// Populate clutch data for types 1v1 to 1v5
		for i := 1; i <= 5; i++ {
			row.AddCell().SetInt(clutchAttempts[i]) // Attempts for clutch type i
			row.AddCell().SetInt(clutchWins[i])     // Wins for clutch type i
		}
	}

	sheet, err = file.AddSheet("Kills")
	checkError(err)

	headerRow = sheet.AddRow()
	headerRow.AddCell().Value = "steamid"
	headerRow.AddCell().Value = "Killer"
	headerRow.AddCell().Value = "Assiter"
	headerRow.AddCell().Value = "Victim"
	headerRow.AddCell().Value = "Weapon"
	headerRow.AddCell().Value = "Distance"
	headerRow.AddCell().Value = "Wallbang"
	headerRow.AddCell().Value = "Headshot"
	headerRow.AddCell().Value = "NoScope"
	headerRow.AddCell().Value = "BlindKill"
	headerRow.AddCell().Value = "Through Smoke"
	headerRow.AddCell().Value = "Assisted Flash"

	for _, kill := range kills {
		assiter := ""
		if kill.Assister != nil {
			assiter = kill.Assister.String()
		}
		row := sheet.AddRow()
		var id = "n/a"
		var killer = "n/a"
		if kill.Weapon.String() != "C4" {
			id = strconv.FormatUint(kill.Killer.SteamID64, 10)
			killer = kill.Killer.String()
		}

		row.AddCell().Value = id
		row.AddCell().SetString(killer)
		row.AddCell().SetString(assiter)
		row.AddCell().SetString(kill.Victim.String())
		row.AddCell().SetString(kill.Weapon.String())
		row.AddCell().SetFloat(float64(kill.Distance))
		row.AddCell().SetBool(kill.IsWallBang())
		row.AddCell().SetBool(kill.IsHeadshot)
		row.AddCell().SetBool(kill.NoScope)
		row.AddCell().SetBool(kill.AttackerBlind)
		row.AddCell().SetBool(kill.ThroughSmoke)
		row.AddCell().SetBool(kill.AssistedFlash)
	}
	currtime := strings.ReplaceAll(time.Now().Format(time.DateTime), " ", "")
	currtime = strings.ReplaceAll(currtime, ":", "-")
	err = file.Save("sheets/" + currtime + ".xlsx")
	checkError(err)
	fmt.Println(Green + "Spreadsheet done " + currtime + ".xlsx" + Reset)

}
