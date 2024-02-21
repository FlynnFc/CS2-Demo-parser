package main

func combineStats(stats []PlayerStats) []PlayerStats {
	mergedStats := make(map[uint64]PlayerStats)

	for _, stats := range stats {
		if existing, found := mergedStats[stats.ID]; found {
			// Merge Basic stats
			existing.Rounds += stats.Rounds
			existing.Matches += stats.Matches
			existing.Kills += stats.Kills
			existing.Deaths += stats.Deaths
			existing.Damage += stats.Damage

			// Merge DetailedStats
			existing.DetailedStats.CTRounds += stats.DetailedStats.CTRounds
			existing.DetailedStats.TRounds += stats.DetailedStats.TRounds
			existing.DetailedStats.RefundTotal += stats.DetailedStats.RefundTotal

			// Merge BombStats
			existing.DetailedStats.BombStats.BombsPlanted += stats.DetailedStats.BombStats.BombsPlanted
			existing.DetailedStats.BombStats.BombsDefused += stats.DetailedStats.BombStats.BombsDefused
			existing.DetailedStats.BombStats.BombDefuseAttempts += stats.DetailedStats.BombStats.BombDefuseAttempts
			existing.DetailedStats.BombStats.BombPlantAttempts += stats.DetailedStats.BombStats.BombPlantAttempts

			// Merge KillStats
			existing.DetailedStats.KillStats.Wallbangs += stats.DetailedStats.KillStats.Wallbangs
			existing.DetailedStats.KillStats.Headshots += stats.DetailedStats.KillStats.Headshots
			existing.DetailedStats.KillStats.Blind += stats.DetailedStats.KillStats.Blind
			existing.DetailedStats.KillStats.NoScope += stats.DetailedStats.KillStats.NoScope
			existing.DetailedStats.KillStats.ThroughSmoke += stats.DetailedStats.KillStats.ThroughSmoke
			existing.DetailedStats.KillStats.AirborneKills += stats.DetailedStats.KillStats.AirborneKills

			// Merge Entry stats
			existing.DetailedStats.KillStats.Entry.CTEntryDeaths += stats.DetailedStats.KillStats.Entry.CTEntryDeaths
			existing.DetailedStats.KillStats.Entry.CTEntryKills += stats.DetailedStats.KillStats.Entry.CTEntryKills
			existing.DetailedStats.KillStats.Entry.TEntryDeaths += stats.DetailedStats.KillStats.Entry.TEntryDeaths
			existing.DetailedStats.KillStats.Entry.TEntryKills += stats.DetailedStats.KillStats.Entry.TEntryKills
			existing.DetailedStats.KillStats.Entry.TotalAttempts += stats.DetailedStats.KillStats.Entry.TotalAttempts

			// Trade Stats
			existing.DetailedStats.KillStats.TradeStats.CTFailedTrades += stats.DetailedStats.KillStats.TradeStats.CTFailedTrades
			existing.DetailedStats.KillStats.TradeStats.CTTradeKills += stats.DetailedStats.KillStats.TradeStats.CTTradeKills
			existing.DetailedStats.KillStats.TradeStats.CTTradedDeaths += stats.DetailedStats.KillStats.TradeStats.CTTradedDeaths
			existing.DetailedStats.KillStats.TradeStats.TFailedTrades += stats.DetailedStats.KillStats.TradeStats.TFailedTrades
			existing.DetailedStats.KillStats.TradeStats.TTradeKills += stats.DetailedStats.KillStats.TradeStats.TTradeKills
			existing.DetailedStats.KillStats.TradeStats.TTradedDeaths += stats.DetailedStats.KillStats.TradeStats.TTradedDeaths
			existing.DetailedStats.KillStats.TradeStats.TotalAttempts += stats.DetailedStats.KillStats.TradeStats.TotalAttempts

			// Merge ClutchStats slice
			existing.DetailedStats.ClutchStats = append(existing.DetailedStats.ClutchStats, stats.DetailedStats.ClutchStats...)

			mergedStats[stats.ID] = existing
		} else {
			// No duplicate, just add to map
			mergedStats[stats.ID] = stats
		}
	}

	// Convert map back to slice
	var mergedSlice []PlayerStats
	for _, stats := range mergedStats {
		mergedSlice = append(mergedSlice, stats)
	}

	return mergedSlice
}
